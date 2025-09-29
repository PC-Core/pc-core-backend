package controllers

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/outputs"
)

func GetRestCharsObject(chars database.ProductChars) (*outputs.RestCharsObject, errors.PCCError) {
	var cc []outputs.RestCharsComponent

	if lc, ok := chars.(*models.LaptopChars); ok {

		cpuChars, cerr := GetRestCharsObject(lc.Cpu)
		gpuChars, gerr := GetRestCharsObject(lc.Gpu)

		if cerr != nil {
			return nil, cerr
		}

		if gerr != nil {
			return nil, gerr
		}

		cc = append(cc, cpuChars.Components...)
		cc = append(cc, gpuChars.Components...)
		cc = append(cc, outputs.RestCharsComponent{Type: "ram", Values: []models.CharsDescription{{Title: "Capacity", Key: "cap"}}, Info: map[string]int16{"cap": lc.Ram}})

		return outputs.NewRestCharsObject(lc.ID, cc), nil
	}

	if gchars, gok := chars.(*models.GpuChars); gok {
		gpuInfo, gerr := GetCharsDescription(gchars)

		if gerr != nil {
			return nil, gerr
		}

		cc = append(cc, outputs.RestCharsComponent{Type: "gpu", Values: gpuInfo, Info: gchars})

		return outputs.NewRestCharsObject(uint64(gchars.ID), cc), nil
	}

	if chars, ok := chars.(*models.CpuChars); ok {
		cpuInfo, err := GetCharsDescription(chars)

		if err != nil {
			return nil, err
		}

		cc = append(cc, outputs.RestCharsComponent{Type: "cpu", Values: cpuInfo, Info: chars})

		return outputs.NewRestCharsObject(chars.ID, cc), nil
	}

	return nil, errors.NewInternalSecretError()
}

func getLaptopCharsDescription() []models.CharsDescription {
	return []models.CharsDescription{
		{Title: "CPU", Key: "cpu"},
		{Title: "RAM", Key: "ram"},
		{Title: "GPU", Key: "gpu"},
	}
}

func getCpuCharsDescription() []models.CharsDescription {
	return []models.CharsDescription{
		{Title: "Name", Key: "name"},
		{Title: "Performance Cores", Key: "pcores"},
		{Title: "Efficiency Cores", Key: "ecores"},
		{Title: "Threads", Key: "threads"},
		{Title: "Base PCores Frequency", Key: "base_p_freq_mhz"},
		{Title: "Max PCores Frequency", Key: "max_p_freq_mhz"},
		{Title: "Base ECores Frequency", Key: "base_e_freq_mhz"},
		{Title: "Max ECores Frequency", Key: "max_e_freq_mhz"},
		{Title: "Socket", Key: "Socket"},
		{Title: "L1 Cache Size", Key: "l1_kb"},
		{Title: "L2 Cache Size", Key: "l2_kb"},
		{Title: "L3 Cache Size", Key: "l3_kb"},
		{Title: "Technical Process", Key: "tecproc_nm"},
		{Title: "TDP", Key: "tdp_watt"},
		{Title: "Release Year", Key: "release_year"},
	}
}

func getGpuCharsDescription() []models.CharsDescription {
	return []models.CharsDescription{
		{Title: "Name", Key: "name"},
		{Title: "Memory", Key: "memory_gb"},
		{Title: "Memory Type", Key: "memory_type"},
		{Title: "Bus Width", Key: "bus_width_bit"},
		{Title: "Base Core Frequency", Key: "base_freq_mhz"},
		{Title: "Boost Core Frequency", Key: "boost_freq_mhz"},
		{Title: "Technical Process", Key: "tecproc_nm"},
		{Title: "TDP", Key: "tdp_watt"},
		{Title: "Release Year", Key: "release_year"},
	}
}

func GetCharsDescription(chars database.ProductChars) ([]models.CharsDescription, errors.PCCError) {
	if _, ok := chars.(*models.LaptopChars); ok {
		return getLaptopCharsDescription(), nil
	}

	if _, ok := chars.(*models.CpuChars); ok {
		return getCpuCharsDescription(), nil
	}

	if _, ok := chars.(*models.GpuChars); ok {
		return getGpuCharsDescription(), nil
	}

	return nil, errors.NewInternalSecretError()
}

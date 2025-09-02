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

		cpuChars, err := GetRestCharsObject(lc.Cpu)

		if err != nil {
			return nil, err
		}

		cc = append(cc, cpuChars.Components...)
		cc = append(cc, outputs.RestCharsComponent{Type: "gpu", Values: []models.CharsDescription{{Title: "Name", Key: "gpu"}}, Info: map[string]string{"gpu": lc.Gpu}})
		cc = append(cc, outputs.RestCharsComponent{Type: "ram", Values: []models.CharsDescription{{Title: "Capacity", Key: "cap"}}, Info: map[string]int16{"cap": lc.Ram}})

		return outputs.NewRestCharsObject(lc.ID, cc), nil
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

func GetCharsDescription(chars database.ProductChars) ([]models.CharsDescription, errors.PCCError) {
	if _, ok := chars.(*models.LaptopChars); ok {
		return []models.CharsDescription{
			{Title: "CPU", Key: "cpu"},
			{Title: "RAM", Key: "ram"},
			{Title: "GPU", Key: "gpu"},
		}, nil
	}

	if _, ok := chars.(*models.CpuChars); ok {
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
		}, nil
	}

	return nil, errors.NewInternalSecretError()
}

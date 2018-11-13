// Copyright 2017-2018 DigitalOcean.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//+build !dragonfly,!freebsd,linux,!netbsd,!openbsd,!solaris,!windows

package cpu

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/digitalocean/go-smbios/smbios"
)

var (
	processorTypeList = map[int]string{
		1: "Other",
		2: "Unknown",
		3: "Central Processor",
		4: "Math Processor",
		5: "DSP Processor",
		6: "Video Processor",
	}

	processorFamilyList = map[int]string{
		1:   "Other",
		2:   "Unknown",
		3:   "8086",
		4:   "80286",
		5:   "Intel386™ processor",
		6:   "Intel486™ processor",
		7:   "8087",
		8:   "80287",
		9:   "80387",
		10:  "80487",
		11:  "Intel® Pentium® processor",
		12:  "Pentium® Pro processor",
		13:  "Pentium® II processor",
		14:  "Pentium® processor with MMX™ technology",
		27:  "K6-2",
		28:  "K6-3",
		29:  "AMD Athlon™ Processor Family [1]",
		30:  "AMD29000 Family",
		31:  "K6-2+",
		32:  "Power PC Family",
		33:  "Power PC 601",
		34:  "Power PC 603",
		35:  "Power PC 603+",
		36:  "Power PC 604",
		37:  "Power PC 620",
		38:  "Power PC x704",
		39:  "Power PC 750",
		40:  "Intel® Core™ Duo processor",
		41:  "Intel® Core™ Duo mobile processor",
		42:  "Intel® Core™ Solo mobile processor",
		43:  "Intel® Atom™ processor",
		44:  "Intel® Core™ M processor",
		45:  "Intel(R) Core(TM) m3 processor",
		46:  "Intel(R) Core(TM) m5 processor",
		47:  "Intel(R) Core(TM) m7 processor",
		48:  "Alpha Family [2]",
		49:  "Alpha 21064",
		50:  "Alpha 21066",
		51:  "Alpha 21164",
		52:  "Alpha 21164PC",
		53:  "Alpha 21164a",
		54:  "Alpha 21264",
		55:  "Alpha 21364",
		56:  "AMD Turion™ II Ultra Dual-Core Mobile M Processor Family",
		57:  "AMD Turion™ II Dual-Core Mobile M Processor Family",
		58:  "AMD Athlon™ II Dual-Core M Processor Family",
		59:  "AMD Opteron™ 6100 Series Processor",
		60:  "AMD Opteron™ 4100 Series Processor",
		61:  "AMD Opteron™ 6200 Series Processor",
		62:  "AMD Opteron™ 4200 Series Processor",
		63:  "AMD FX™ Series Processor",
		64:  "MIPS Family",
		65:  "MIPS R4000",
		66:  "MIPS R4200",
		67:  "MIPS R4400",
		68:  "MIPS R4600",
		69:  "MIPS R10000",
		70:  "AMD C-Series Processor",
		71:  "AMD E-Series Processor",
		72:  "AMD A-Series Processor",
		73:  "AMD G-Series Processor",
		74:  "AMD Z-Series Processor",
		75:  "AMD R-Series Processor",
		76:  "AMD Opteron™ 4300 Series Processor",
		77:  "AMD Opteron™ 6300 Series Processor",
		78:  "AMD Opteron™ 3300 Series Processor",
		79:  "AMD FirePro™ Series Processor",
		80:  "SPARC Family",
		81:  "SuperSPARC",
		82:  "microSPARC II",
		83:  "microSPARC IIep",
		84:  "UltraSPARC",
		85:  "UltraSPARC II",
		86:  "UltraSPARC Iii",
		87:  "UltraSPARC III",
		88:  "UltraSPARC IIIi",
		96:  "68040 Family",
		97:  "68xxx",
		98:  "68000",
		99:  "68010",
		100: "68020",
		101: "68030",
		102: "AMD Athlon(TM) X4 Quad-Core Processor Family",
		103: "AMD Opteron(TM) X1000 Series Processor",
		104: "AMD Opteron(TM) X2000 Series APU",
		105: "AMD Opteron(TM) A-Series Processor",
		106: "AMD Opteron(TM) X3000 Series APU",
		107: "AMD Zen Processor Family",
		112: "Hobbit Family",
		120: "Crusoe™ TM5000 Family",
		121: "Crusoe™ TM3000 Family",
		122: "Efficeon™ TM8000 Family",
		128: "Weitek",
		130: "Itanium™ processor",
		131: "AMD Athlon™ 64 Processor Family",
		132: "AMD Opteron™ Processor Family",
		133: "AMD Sempron™ Processor Family",
		134: "AMD Turion™ 64 Mobile Technology",
		135: "Dual-Core AMD Opteron™ Processor Family",
		136: "AMD Athlon™ 64 X2 Dual-Core Processor Family",
		137: "AMD Turion™ 64 X2 Mobile Technology",
		138: "Quad-Core AMD Opteron™ Processor Family",
		139: "Third-Generation AMD Opteron™ Processor Family",
		140: "AMD Phenom™ FX Quad-Core Processor Family",
		141: "AMD Phenom™ X4 Quad-Core Processor Family",
		142: "AMD Phenom™ X2 Dual-Core Processor Family",
		143: "AMD Athlon™ X2 Dual-Core Processor Family",
		144: "PA-RISC Family",
		145: "PA-RISC 8500",
		146: "PA-RISC 8000",
		147: "PA-RISC 7300LC",
		148: "PA-RISC 7200",
		149: "PA-RISC 7100LC",
		150: "PA-RISC 7100",
		160: "V30 Family",
		161: "Quad-Core Intel® Xeon® processor 3200 Series",
		162: "Dual-Core Intel® Xeon® processor 3000 Series",
		163: "Quad-Core Intel® Xeon® processor 5300 Series",
		164: "Dual-Core Intel® Xeon® processor 5100 Series",
		165: "Dual-Core Intel® Xeon® processor 5000 Series",
		166: "Dual-Core Intel® Xeon® processor LV",
		167: "Dual-Core Intel® Xeon® processor ULV",
		168: "Dual-Core Intel® Xeon® processor 7100 Series",
		169: "Quad-Core Intel® Xeon® processor 5400 Series",
		170: "Quad-Core Intel® Xeon® processor",
		171: "Dual-Core Intel® Xeon® processor 5200 Series",
		172: "Dual-Core Intel® Xeon® processor 7200 Series",
		173: "Quad-Core Intel® Xeon® processor 7300 Series",
		174: "Quad-Core Intel® Xeon® processor 7400 Series",
		177: "Pentium® III Processor with Intel® SpeedStep™ Technology",
		178: "Pentium® 4 Processor",
		179: "Intel® Xeon® processor",
		180: "AS400 Family",
		181: "Intel® Xeon™ processor MP",
		182: "AMD Athlon™ XP Processor Family",
		183: "AMD Athlon™ MP Processor Family",
		184: "Intel® Itanium® 2 processor",
		185: "Intel® Pentium® M processor",
		186: "Intel® Celeron® D processor",
		187: "Intel® Pentium® D processor",
		188: "Intel® Pentium® Processor Extreme Edition",
		189: "Intel® Core™ Solo Processor",
		191: "Intel® Core™ 2 Duo Processor",
		192: "Intel® Core™ 2 Solo processor",
		193: "Intel® Core™ 2 Extreme processor",
		194: "Intel® Core™ 2 Quad processor",
		195: "Intel® Core™ 2 Extreme mobile processor",
		196: "Intel® Core™ 2 Duo mobile processor",
		197: "Intel® Core™ 2 Solo mobile processor",
		198: "Intel® Core™ i7 processor",
		199: "Dual-Core Intel® Celeron® processor",
		200: "IBM390 Family",
		201: "G4",
		202: "G5",
		203: "ESA/390 G6",
		204: "z/Architecture base",
		205: "Intel® Core™ i5 processor",
		206: "Intel® Core™ i3 processor",
		210: "VIA C7™-M Processor Family",
		211: "VIA C7™-D Processor Family",
		212: "VIA C7™ Processor Family",
		213: "VIA Eden™ Processor Family",
		214: "Multi-Core Intel® Xeon® processor",
		215: "Dual-Core Intel® Xeon® processor 3xxx Series",
		216: "Quad-Core Intel® Xeon® processor 3xxx Series",
		217: "VIA Nano™ Processor Family",
		218: "Dual-Core Intel® Xeon® processor 5xxx Series",
		219: "Quad-Core Intel® Xeon® processor 5xxx Series",
		221: "Dual-Core Intel® Xeon® processor 7xxx Series",
		222: "Quad-Core Intel® Xeon® processor 7xxx Series",
		223: "Multi-Core Intel® Xeon® processor 7xxx Series",
		224: "Multi-Core Intel® Xeon® processor 3400 Series",
		228: "AMD Opteron™ 3000 Series Processor",
		229: "AMD Sempron™ II Processor",
		230: "Embedded AMD Opteron™ Quad-Core Processor Family",
		231: "AMD Phenom™ Triple-Core Processor Family",
		232: "AMD Turion™ Ultra Dual-Core Mobile Processor Family",
		233: "AMD Turion™ Dual-Core Mobile Processor Family",
		234: "AMD Athlon™ Dual-Core Processor Family",
		235: "AMD Sempron™ SI Processor Family",
		236: "AMD Phenom™ II Processor Family",
		237: "AMD Athlon™ II Processor Family",
		238: "Six-Core AMD Opteron™ Processor Family",
		239: "AMD Sempron™ M Processor Family",
		250: "i860",
		251: "i960",
		256: "ARMv7",
		257: "ARMv8",
		260: "SH-3",
		261: "SH-4",
		280: "ARM",
		281: "StrongARM",
		300: "6x86",
		301: "MediaGX",
		302: "MII",
		320: "WinChip",
		350: "DSP",
		500: "Video Processor",
	}

	processorUpgradeList = map[int]string{
		1:  "Other",
		2:  "Unknown",
		3:  "Daughter Board",
		4:  "ZIF Socket",
		5:  "Replaceable Piggy Back",
		6:  "None",
		7:  "LIF Socket",
		8:  "Slot 1",
		9:  "Slot 2",
		10: "370-pin socket",
		11: "Slot A",
		12: "Slot M",
		13: "Socket 423",
		14: "Socket A (Socket 462)",
		15: "Socket 478",
		16: "Socket 754",
		17: "Socket 940",
		18: "Socket 939",
		19: "Socket mPGA604",
		20: "Socket LGA771",
		21: "Socket LGA775",
		22: "Socket S1",
		23: "Socket AM2",
		24: "Socket F (1207)",
		25: "Socket LGA1366",
		26: "Socket G34",
		27: "Socket AM3",
		28: "Socket C32",
		29: "Socket LGA1156",
		30: "Socket LGA1567",
		31: "Socket PGA988A",
		32: "Socket BGA1288",
		33: "Socket rPGA988B",
		34: "Socket BGA1023",
		35: "Socket BGA1224",
		36: "Socket LGA1155",
		37: "Socket LGA1356",
		38: "Socket LGA2011",
		39: "Socket FS1",
		40: "Socket FS2",
		41: "Socket FM1",
		42: "Socket FM2",
		43: "Socket LGA2011-3",
		44: "Socket LGA1356-3",
		45: "Socket LGA1150",
		46: "Socket BGA1168",
		47: "Socket BGA1234",
		48: "Socket BGA1364",
		49: "Socket AM4",
		50: "Socket LGA1151",
		51: "Socket BGA1356",
		52: "Socket BGA1440",
		53: "Socket BGA1515",
		54: "Socket LGA3647-1",
		55: "Socket SP3",
		56: "Socket SP3r2",
	}
	pcDescList = map[string]string{
		"PC_RSVD": "Reserved",
		"PC_UNK":  "Unknown",
		"PC_64B":  "64-bit Capable",
		"PC_MC":   "Multi-Core",
		"PC_HT":   "Hardware Thread",
		"PC_EP":   "Execute Protection",
		"PC_EV":   "Enhanced Virtualization",
		"PC_PPC":  "Power/Performance Control",
	}
	pcIntList = map[string]int{
		"PC_RSVD": 1,
		"PC_UNK":  2,
		"PC_64B":  4,
		"PC_MC":   8,
		"PC_HT":   16,
		"PC_EP":   32,
		"PC_EV":   64,
		"PC_PPC":  128,
	}
	pfFlagList = map[int]string{
		0:  "fpu",
		1:  "vme",
		2:  "de",
		3:  "pse",
		4:  "tsc",
		5:  "msr",
		6:  "pae",
		7:  "mce",
		8:  "cx8",
		9:  "apic",
		11: "sep",
		12: "mtrr",
		13: "pge",
		14: "mca",
		15: "cmov",
		16: "pat",
		17: "pse-36",
		18: "psn",
		19: "clfsh",
		21: "ds",
		22: "acpi",
		23: "mmx",
		24: "fxsr",
		25: "sse",
		26: "sse2",
		27: "ss",
		28: "htt",
		29: "tm",
		30: "ia64",
		31: "pbe",
	}

	cpuStatusFlags = map[int]string{
		0: "Unknown",
		1: "CPU Enabled",
		2: "CPU Disabled by User through BIOS Setup",
		3: "CPU Disabled by BIOS (POST Error)",
		4: "CPU is idle, waiting to be enabled",
		7: "Other",
	}
)

func getCPUFeatureFlags(val int) (flags []string) {
	for key := range pfFlagList {
		if (byte(val) & byte(key)) != 0 {
			flags = append(flags, pfFlagList[key])
		}
	}
	return flags
}

// CPU Characteristics flags are defined by 2 bytes worth of bits (0-15)
// As at v3.1.1 of SMBIOS, bits 8 to 15 are 'reserved'
func getCPUCharacteristicsFlags(val int) (flags []string) {
	for key := range pcIntList {
		if (byte(val) & byte(pcIntList[key])) != 0 {
			flags = append(flags, pcDescList[key])
		}
	}
	return flags
}

func getCPUStatusFlags(val int) (flags []string) {
	if (val & 64) != 0 {
		flags = append(flags, "Socket populated")
	} else {
		flags = append(flags, "Socket unpopulated")
	}
	flags = append(flags, cpuStatusFlags[int(val&^240)])
	return flags
}

func getCPUVoltage(val int) (voltage float32) {
	if byte(val)&0x80 != 0 {
		//MSB is 1 therefore voltage is ((val - 128)/10) volts
		voltage = ((float32(val) - 128) / 10)
	} else {
		//MSB is 0 and therefore voltage is defined by bit flag
		if byte(val)&0x01 != 0 {
			voltage = 5.0
		} else if byte(val)&0x02 != 0 {
			voltage = 3.3
		} else if byte(val)&0x03 != 0 {
			voltage = 2.9
		}
	}
	return voltage
}

// The CPU Model Int uses two half bytes:
// First 4 bits of s.Formatted[4]
// Last 4 bits of s.Formatted[6]
func getCPUModelInt(val1, val2 int) int {
	return int((((val1 &^ 15) >> 4) + ((val2 &^ 160) << 4)))
}

// Get Function to build a *Cpu struct object with all
// the details from SMBIOS
func (cpu *CPU) Get(s *smbios.Structure) error {
	// SMBIOS returns an index starting from 1, need to -1 for Go slice indices
	cpu.SocketDesignation = strings.TrimSpace(s.Strings[(s.Formatted[0] - 1)])
	cpu.ProcessorManufacturer = strings.TrimSpace(s.Strings[(s.Formatted[3] - 1)])
	cpu.Version = strings.TrimSpace(s.Strings[(s.Formatted[12] - 1)])
	cpu.SerialNumber = strings.TrimSpace(s.Strings[(s.Formatted[28] - 1)])
	cpu.AssetTag = strings.TrimSpace(s.Strings[(s.Formatted[29] - 1)])
	cpu.PartNumber = strings.TrimSpace(s.Strings[(s.Formatted[30] - 1)])
	cpu.ProcessorType = processorTypeList[int(s.Formatted[1])]
	cpu.ProcessorFamily = processorFamilyList[int(s.Formatted[2])]
	cpu.Stepping = int(s.Formatted[4]) & 15
	cpu.Family = int(s.Formatted[5]) & 15
	cpu.Model = getCPUModelInt(int(s.Formatted[4]), int(s.Formatted[6]))
	cpu.Voltage = getCPUVoltage(int(s.Formatted[13]))
	cpu.ExternalClock = int(binary.LittleEndian.Uint16(s.Formatted[14:16]))
	cpu.MaxSpeed = int(binary.LittleEndian.Uint16(s.Formatted[16:18]))
	cpu.CurrentSpeed = int(binary.LittleEndian.Uint16(s.Formatted[18:20]))
	cpu.L1CacheHandle = fmt.Sprintf("0x%04d", binary.LittleEndian.Uint16(s.Formatted[22:24]))
	cpu.L2CacheHandle = fmt.Sprintf("0x%04d", binary.LittleEndian.Uint16(s.Formatted[24:26]))
	cpu.L3CacheHandle = fmt.Sprintf("0x%04d", binary.LittleEndian.Uint16(s.Formatted[26:28]))
	cpu.StatusFlags = getCPUStatusFlags(int(s.Formatted[20]))
	cpu.ProcessorUpgrade = processorUpgradeList[int(s.Formatted[21])]
	cpu.CoreCount = int(s.Formatted[31])
	cpu.CoreEnabled = int(s.Formatted[32])
	cpu.ThreadCount = int(s.Formatted[33])
	// Fields 34 and 35 make up the bits for the characteristics flags
	cpu.ProcessorCharacteristics = getCPUCharacteristicsFlags(int(binary.LittleEndian.Uint16(s.Formatted[34:36])))
	// Fields 8 - 11 make up the bits for the EDX CPU Feature flags
	cpu.ProcessorFlags = getCPUFeatureFlags(int(binary.LittleEndian.Uint16(s.Formatted[8:12])))

	return nil
}

// CPU Structure for containing Processor information
type CPU struct {
	SocketDesignation        string
	ProcessorType            string
	ProcessorFamily          string
	ProcessorManufacturer    string
	Stepping                 int
	Model                    int
	Family                   int
	Type                     int
	Version                  string
	Voltage                  float32
	ExternalClock            int
	MaxSpeed                 int
	CurrentSpeed             int
	StatusFlags              []string
	ProcessorUpgrade         string
	L1CacheHandle            string
	L2CacheHandle            string
	L3CacheHandle            string
	SerialNumber             string
	AssetTag                 string
	PartNumber               string
	CoreCount                int
	CoreEnabled              int
	ThreadCount              int
	ProcessorFlags           []string
	ProcessorCharacteristics []string
}

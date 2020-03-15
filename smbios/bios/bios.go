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


package bios

import (
        "encoding/binary"
        "fmt"
        "strings"

        "github.com/axrayn/go-smbios/smbios"
)

type Bios struct {                                                                                                              
	    Vendor                  string
        Version                 string
        StartingAddressSegment  string
        ReleaseDate             string
        ROMSize                 int
        Characteristics         []string
        CharacteristicsExtended []string
        MajorRelease            int
        MinorRelease            int
        FirmwareMajorRelease    int
        FirmwareMinorRelease    int
        ROMSizeExtended         int
}

var (
	biosCharacterList = map[int]string{
	  0: "Reserved",
	  1: "Reserved",
	  2: "Unknown",
	  3: "Not Supported",
	  4: "ISA",
	  5: "MCA",
	  6: "EISA",
	  7: "PCI",
	  8: "PCMCIA",
	  9: "PnP",
	  10: "APM",
	  11: "Upgradeable (Flash)",
	  12: "Shadowing",
	  13: "VL-VESA",
	  14: "ESCD",
	  15: "CD Boot",
	  16: "Selectable Boot",
	  17: "BIOS ROM Socketed",
	  18: "PCMCIA Boot",
	  19: "EDD",
	  20: "NEC Floppy",
	  21: "Toshiba Floppy",
	  22: "5.25in/360KB Floppy",
	  23: "5.25in/1.2MB Floppy",
	  24: "3.5in/720KB Floppy",
	  25: "3.5in/2.88MB Floppy",
	  26: "PrintScreen Service",
	  27: "8042 Keyboard",
	  28: "Serial Services",
	  29: "Printer Services",
	  30: "CGA/Mono Video Services",
	  31: "NEC PC-98",
	}
	biosCharacterEx1List = map[int]string{
	  0: "ACPI",
	  1: "USB Legacy",
	  2: "AGP",
	  3: "I2O Boot",
	  4: "LS-120 SuperDisk Boot",
	  5: "ATAPI ZIP Drive Boot",
	  6: "1394 Boot",
	  7: "Smart Battery",
	}
	biosCharacterExList2 = map[int]string{
	  0: "BIOS Boot Specification",
	  1: "Fn Key Network Boot",
	  2: "Targetted Content Distribution",
	  3: "UEFI Specification",
	  4: "IsVirtual",
	}
)

func getBIOSCharacteristicsFlags(val int) (flags []string) {
	for key := range biosCharacterList {
			if (byte(val) & byte(key)) != 0 {
					flags = append(flags, biosCharacterList[key])
			}
	}
	return flags
}


// Get Function to build a *Bios struct object with all
// the details from SMBIOS
func (bios *Bios) Get(s *smbios.Structure) error {
	// SMBIOS returns an index starting from 1, need to -1 for Go slice indices
	bios.Vendor = strings.TrimSpace(s.Strings[s.Formatted[0] - 1])
	bios.Version = strings.TrimSpace(s.Strings[s.Formatted[1] - 1])
	bios.ReleaseDate = strings.TrimSpace(s.Strings[s.Formatted[4] - 1])
	bios.StartingAddressSegment = fmt.Sprintf("Starting Address: 0x%04X\n", int(binary.LittleEndian.Uint16(s.Formatted[2:4])))
	// ROM Size is either here or in the extended bit, depending on the value here being FFh or not
	if (s.Formatted[5] > 254) {
		bios.ROMSize = int(s.Formatted[5] + 1)*64
	} else {
			bios.ROMSize = int(s.Formatted[5] + 1)*64
	}
	bios.Characteristics = getBIOSCharacteristicsFlags(int(binary.LittleEndian.Uint16(s.Formatted[14:16])))
	bios.MajorRelease            = int(s.Formatted[16])
	bios.MinorRelease            = int(s.Formatted[17])
	if (int(s.Formatted[18]) < 255) {
			bios.FirmwareMajorRelease    = int(s.Formatted[18])
			bios.FirmwareMinorRelease    = int(s.Formatted[19])
	}
	//bios.CharacteristicsExtended = getBIOSCharacteristicsFlags(int(binary.LittleEndian.Uint16(s.Formatted[14:16])))
	//bios.ROMSizeExtended         byte
	return nil
}

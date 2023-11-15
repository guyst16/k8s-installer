package virtualmachine

import (
	"errors"
	"fmt"
	"log"

	"github.com/digitalocean/go-libvirt"
	"github.com/guyst16/mykube/pkg/libvirtconn"
	"golang.org/x/crypto/ssh"
)

type Virtualmachine struct {
	os_name          string
	os_path          string
	cloudconfig_path string
	vcpu_amount      int
	memory_amount    int
	name             string
}

var STATES = map[libvirt.ConnectListAllDomainsFlags]string{libvirt.ConnectListDomainsRunning: "Running", libvirt.ConnectListDomainsPaused: "Paused", libvirt.ConnectListDomainsShutoff: "Shutoff"}

// Create virtual machine object
func NewVirtualmachine(os_name string, os_path string, cloudconfig_path string, vcpu_amount int, memory_amount int, name string) *Virtualmachine {
	v := Virtualmachine{os_name: os_name, os_path: os_path, cloudconfig_path: cloudconfig_path, vcpu_amount: vcpu_amount, memory_amount: memory_amount, name: name}
	return &v
}

func GetVirtualMachine(vmName string) (dom *libvirt.Domain) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	vm, _ := libvirtconn.DomainLookupByName(vmName)
	return &vm
}

func ListAllVirtualmachines() {
	println("ID\tNAME\t\tUUID\t\t\t\t\tSTATE")
	println("-----------------------------------------------------------------------")
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	for state := range STATES {
		flags := state
		domains, _, _ := libvirtconn.ConnectListAllDomains(1, flags)
		for _, vm := range domains {
			fmt.Printf("%d\t%s\t%x\t%s\n", vm.ID, vm.Name, vm.UUID, STATES[state])
		}
	}
}

func (vm Virtualmachine) CreateVirtualmachine() {
	vmXML := ModifyXML("assets/vmTemplate.xml", vm.name, vm.os_path, vm.cloudconfig_path)
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	vmXMLString := string(vmXML)
	_, err := libvirtconn.DomainDefineXML(vmXMLString)
	if err != nil {
		log.Fatal(err)
	}
}

// Start defined vm
func StartVirtualMachine(vmName string) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	domain := GetVirtualMachine(vmName)
	if domain == nil {
		log.Fatal("Virtual machine not defined")
	}
	err := libvirtconn.DomainCreate(*domain)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete defined vm
func DeleteVirtualMachine(vmName string) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	domain := GetVirtualMachine(vmName)
	if domain == nil {
		log.Fatal("Virtual machine not defined")
	}
	err := libvirtconn.DomainUndefine(*domain)
	if err != nil {
		log.Fatal(err)
	}
	err = libvirtconn.DomainDestroy(*domain)
	if err != nil {
		log.Fatal(err)
	}
}

// Get virtual machine public IP
func GetVirtualMachineIP(vmName string) (vmIPAddress string, err error) {
	libvirtconn := libvirtconn.ConnectLibvirtLocal()
	domain := GetVirtualMachine(vmName)
	domainIPAddress := ""
	if domain == nil {
		return domainIPAddress, errors.New("virtual machine not defined")
	}

	// Get list of the domain interfaces addresses
	intrefacesList, err := libvirtconn.DomainInterfaceAddresses(*domain, 0, 0)
	if err != nil {
		return domainIPAddress, err
	}

	if len(intrefacesList) == 0 {
		return domainIPAddress, errors.New("domain " + vmName + " doesn't have an IP address")
	}

	return intrefacesList[0].Addrs[0].Addr, nil
}

// Get ssh client for a virtaul machine
func GetVirtualMachineSSHConnection(vmName string) {
	var hostKey ssh.PublicKey

	config := &ssh.ClientConfig{
		User: "sumit",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(""),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	client, err := ssh.Dial("tcp", "yourserver.com:22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

}

// Create ssh key for a virtual machine
func CreateVirtualmachineSHHKey(vmName string, outputKeyPath string)

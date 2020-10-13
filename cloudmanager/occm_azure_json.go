package cloudmanager

func (c *Client) callParameters() string {
	return `{
        "location": {
            "value": "string"
        },
        "virtualMachineName": {
            "value": "string"
        },
        "adminUsername": {
            "value": "string"
        },
        "virtualNetworkId": {
            "value": "string"
        },
        "adminPassword": {
            "value": "string"
        },
        "subnetId": {
            "value": "string"
        },
        "customData": {
        "value": "string"
        },
        "environment": {
            "value": "string"
        }
    }`
}

func (c *Client) callTemplate() string {
	return `{
        "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
        "contentVersion": "1.0.0.0",
        "parameters": {
            "location": {
                "type": "string",
                "defaultValue": "eastus"
            },
            "virtualMachineName": {
                "type": "string"
            },
            "adminUsername": {
                "type": "string"
            },
            "virtualNetworkId": {
                "type": "string"
            },
            "networkSecurityGroupName": {
                "type": "string",
                "defaultValue": "[concat(parameters('virtualMachineName'),'-nsg')]"
            },
            "adminPassword": {
                "type": "securestring"
            },
            "subnetId": {
                "type": "string"
            },
            "customData": {
              "type": "string"
            },
            "environment": {
              "type": "string",
              "defaultValue": "prod"
            }
        },
        "variables": {
            "vnetId": "[parameters('virtualNetworkId')]",
            "subnetRef": "[parameters('subnetId')]",
            "networkInterfaceName": "[concat(parameters('virtualMachineName'),'-nic')]",
            "diagnosticsStorageAccountName": "[concat(toLower(parameters('virtualMachineName')),'sa')]",
            "diagnosticsStorageAccountId": "[concat('Microsoft.Storage/storageAccounts/', variables('diagnosticsStorageAccountName'))]",
            "diagnosticsStorageAccountType": "Standard_LRS",
            "publicIpAddressName": "[concat(parameters('virtualMachineName'),'-ip')]",
            "publicIpAddressType": "Dynamic",
            "publicIpAddressSku": "Basic",
            "msiExtensionName": "ManagedIdentityExtensionForLinux",
            "virtualMachineSize": "Standard_D2s_v3",
            "occmNSG": "[concat(parameters('virtualMachineName'),'-nsg')]",
            "occmOffer": "[if(equals(parameters('environment'), 'stage'), 'netapp-oncommand-cloud-manager-staging-preview', 'netapp-oncommand-cloud-manager')]"
        },
        "resources": [
            {
                "name": "[parameters('virtualMachineName')]",
                "type": "Microsoft.Compute/virtualMachines",
                "apiVersion": "2018-04-01",
                "location": "[parameters('location')]",
                "dependsOn": [
                    "[concat('Microsoft.Network/networkInterfaces/', variables('networkInterfaceName'))]",
                    "[concat('Microsoft.Storage/storageAccounts/', variables('diagnosticsStorageAccountName'))]"
                ],
                "properties": {
                    "osProfile": {
                        "computerName": "[parameters('virtualMachineName')]",
                        "adminUsername": "[parameters('adminUsername')]",
                        "adminPassword": "[parameters('adminPassword')]",
                        "customData": "[base64(parameters('customData'))]"
                    },
                    "hardwareProfile": {
                        "vmSize": "[variables('virtualMachineSize')]"
                    },
                    "storageProfile": {
                        "imageReference": {
                            "publisher": "netapp",
                            "offer": "[variables('occmOffer')]",
                            "sku": "occm-byol",
                            "version": "latest"
                        },
                        "osDisk": {
                            "createOption": "fromImage",
                            "managedDisk": {
                                "storageAccountType": "Premium_LRS"
                            }
                        },
                        "dataDisks": []
                    },
                    "networkProfile": {
                        "networkInterfaces": [
                            {
                                "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('networkInterfaceName'))]"
                            }
                        ]
                    },
                    "diagnosticsProfile": {
                      "bootDiagnostics": {
                        "enabled": true,
                        "storageUri":
                          "[concat('https://', variables('diagnosticsStorageAccountName'), '.blob.core.windows.net/')]"
                      }
                    }
                },
                "plan": {
                    "name": "occm-byol",
                    "publisher": "netapp",
                    "product": "[variables('occmOffer')]"
                },
                "identity": {
                    "type": "systemAssigned"
                }
            },
            {
                "apiVersion": "2017-12-01",
                "type": "Microsoft.Compute/virtualMachines/extensions",
                "name": "[concat(parameters('virtualMachineName'),'/', variables('msiExtensionName'))]",
                "location": "[parameters('location')]",
                "dependsOn": [
                    "[concat('Microsoft.Compute/virtualMachines/', parameters('virtualMachineName'))]"
                ],
                "properties": {
                    "publisher": "Microsoft.ManagedIdentity",
                    "type": "[variables('msiExtensionName')]",
                    "typeHandlerVersion": "1.0",
                    "autoUpgradeMinorVersion": true,
                    "settings": {
                        "port": 50342
                    }
                }
            },
            {
                "name": "[variables('diagnosticsStorageAccountName')]",
                "type": "Microsoft.Storage/storageAccounts",
                "apiVersion": "2015-06-15",
                "location": "[parameters('location')]",
                "properties": {
                  "accountType": "[variables('diagnosticsStorageAccountType')]"
                }
            },
            {
                "name": "[variables('networkInterfaceName')]",
                "type": "Microsoft.Network/networkInterfaces",
                "apiVersion": "2018-04-01",
                "location": "[parameters('location')]",
                "dependsOn": [
                    "[concat('Microsoft.Network/publicIpAddresses/', variables('publicIpAddressName'))]"
                ],
                "properties": {
                    "ipConfigurations": [
                        {
                            "name": "ipconfig1",
                            "properties": {
                                "subnet": {
                                    "id": "[variables('subnetRef')]"
                                },
                                "privateIPAllocationMethod": "Dynamic",
                                "publicIpAddress": {
                                    "id": "[resourceId(resourceGroup().name,'Microsoft.Network/publicIpAddresses', variables('publicIpAddressName'))]"
                                }
                            }
                        }
                    ],
                    "networkSecurityGroup": {
                        "id": "[resourceId(resourceGroup().name, 'Microsoft.Network/networkSecurityGroups', parameters('networkSecurityGroupName'))]"
                    }
                }
            },
            {
                "name": "[variables('publicIpAddressName')]",
                "type": "Microsoft.Network/publicIpAddresses",
                "apiVersion": "2017-08-01",
                "location": "[parameters('location')]",
                "properties": {
                    "publicIpAllocationMethod": "[variables('publicIpAddressType')]"
                },
                "sku": {
                    "name": "[variables('publicIpAddressSku')]"
                }
            },
            {
                "name": "[parameters('networkSecurityGroupName')]",
                "condition": "[equals(parameters('networkSecurityGroupName'),variables('occmNSG'))]",
                "type": "Microsoft.Network/networkSecurityGroups",
                "apiVersion": "2018-01-01",
                "location": "[parameters('location')]",
                "properties": {
                    "securityRules": [
                        {
                            "name": "http",
                            "properties": {
                                "priority": 1010,
                                "protocol": "TCP",
                                "access": "Allow",
                                "direction": "Inbound",
                                "sourceApplicationSecurityGroups": [],
                                "destinationApplicationSecurityGroups": [],
                                "sourceAddressPrefix": "*",
                                "sourcePortRange": "*",
                                "destinationAddressPrefix": "*",
                                "destinationPortRange": "80"
                            }
                        },
                        {
                            "name": "https",
                            "properties": {
                                "priority": 1020,
                                "protocol": "TCP",
                                "access": "Allow",
                                "direction": "Inbound",
                                "sourceApplicationSecurityGroups": [],
                                "destinationApplicationSecurityGroups": [],
                                "sourceAddressPrefix": "*",
                                "sourcePortRange": "*",
                                "destinationAddressPrefix": "*",
                                "destinationPortRange": "443"
                            }
                        },
                        {
                            "name": "ssh",
                            "properties": {
                                "priority": 1030,
                                "protocol": "TCP",
                                "access": "Allow",
                                "direction": "Inbound",
                                "sourceApplicationSecurityGroups": [],
                                "destinationApplicationSecurityGroups": [],
                                "sourceAddressPrefix": "*",
                                "sourcePortRange": "*",
                                "destinationAddressPrefix": "*",
                                "destinationPortRange": "22"
                            }
                        }
                    ]
                }
            }
        ],
        "outputs": {
            "publicIpAddressName": {
                "type": "string",
                "value": "[variables('publicIpAddressName')]"
            }
        }
    }`
}

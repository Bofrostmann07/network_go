# network_go

## Query language
### Switch query keys
| Key       | Description                              | Version |
|-----------|------------------------------------------|---------|
| address   | IPv4, IPv6 or FQDN address of the switch | 1.0.0   |
| hostname  | Hostname of the switch                   | 1.0.0   |
| platform  | cisco_ios                                | 1.0.0   |
| group     | Group the switch is a member of          | 1.0.0   |
| reachable | Is switch SSH reachable                  | 1.0.0   |

### Interface filter keys
| Key     | Description                                               | Version |
|---------|-----------------------------------------------------------|---------|
| config  | Search for interface config like 'switchport mode access' | 1.0.0   |
| maclist |                                                           | N/A     |

### Comparison operators
| Operator | Description      | Version |
|----------|------------------|---------|
| =        | Is Equal         | 1.0.0   |
| !=       | Is not Equal     | 1.0.0   |
| ~        | Contains         | 1.0.0   |
| !~       | Does not contain | 1.0.0   |

### Logical operators
| Operator | Description           | Version |
|----------|-----------------------|---------|
| &        | Logical AND           | 1.0.0   |
| &#124;   | Logical OR            | 1.0.0   |
| (...)    | Group querys together | 1.0.0   |

### Examples
#### Switch query
- hostname ~ ".building1.de.company.de" ➔ Lists all switches, which contain the query. E.g. 'switch1.building1.de.company.de' is matching
- hostname ~ ".building1.de.company.de" & group = "core" ➔ Lists all switches, which contain the hostname query and are member of the core group.
- (address ~ "10.10" | address ~ "10.40") & group = "production" ➔ Lists all switches, which start with the IPv4 address '10.10' or '10.40' and are member of the production group.

#### Interface filter
- config = "dot1x pae authenticator" ➔ Lists all interfaces, which have AAA 802.1x activated
- config ~ "switchport access vlan 10" | config ~ "switchport access vlan 20" ➔ Lists all interfaces, which have access vlan 10 or 20 configured
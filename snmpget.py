"""
| $ snmpget -v1 -c public demo.snmplabs.com SNMPv2-MIB::sysDescr.0
"""#
from pysnmp.hlapi import *
import sys

ip,community=sys.argv[1],sys.argv[2]
oID=sys.argv[3]

errorIndication, errorStatus, errorIndex, varBinds = next(
    getCmd(SnmpEngine(),
           CommunityData(community),
           UdpTransportTarget((ip, 161)),
           ContextData(),
           ObjectType(ObjectIdentity(oID)))
)

if errorIndication:
    print(errorIndication)
elif errorStatus:
    print('%s at %s' % (errorStatus.prettyPrint(),
                        errorIndex and varBinds[int(errorIndex) - 1][0] or '?'))
else:
    for varBind in varBinds:
            varB=(' = '.join([x.prettyPrint() for x in varBind]))
            print varB

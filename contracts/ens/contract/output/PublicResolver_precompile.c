#include "./vntlib.h"

KEY string INTERFACE_META_ID = "0x01ffc9a7";
KEY string ADDR_INTERFACE_ID = "0x3b3b57de";
KEY string CONTENT_INTERFACE_ID = "0xd8389dc5";
KEY string NAME_INTERFACE_ID = "0x691f3431";
KEY string ABI_INTERFACE_ID = "0x2203ab56";
KEY string PUBKEY_INTERFACE_ID = "0xc8690233";
KEY string TEXT_INTERFACE_ID = "0x59d1d43c";

EVENT AddrChanged(indexed string node, address a);
EVENT ContentChanged(indexed string node, string hash);
EVENT NameChanged(indexed string node, string name);
EVENT ABIChanged(indexed string node, indexed uint256 contentType);
EVENT PubkeyChanged(indexed string node, string x, string y);
EVENT TextChanged(indexed string node, indexed string indexedKey, string key);

typedef struct
{
    string x;
    string y;
} PublicKey;

typedef struct
{
    address addr;
    string content;
    string name;
    PublicKey pubkey;
    mapping(string, string) text;
    mapping(uint256, string) abis;
} Record;

KEY address ens;
KEY mapping(string, Record) records;

CALL address owner(CallParams params, string node);

   
               
                                             
   

void keyvbvhiap7(){
AddKeyInfo( &records.value.addr, 7, &records, 9, false);
AddKeyInfo( &records.value.addr, 7, &records.key, 6, false);
AddKeyInfo( &records.value.addr, 7, &records.value.addr, 9, false);
AddKeyInfo( &records.value.content, 6, &records, 9, false);
AddKeyInfo( &records.value.content, 6, &records.key, 6, false);
AddKeyInfo( &records.value.content, 6, &records.value.content, 9, false);
AddKeyInfo( &records.value.pubkey.x, 6, &records, 9, false);
AddKeyInfo( &records.value.pubkey.x, 6, &records.key, 6, false);
AddKeyInfo( &records.value.pubkey.x, 6, &records.value.pubkey, 9, false);
AddKeyInfo( &records.value.pubkey.x, 6, &records.value.pubkey.x, 9, false);
AddKeyInfo( &CONTENT_INTERFACE_ID, 6, &CONTENT_INTERFACE_ID, 9, false);
AddKeyInfo( &records.value.abis.value, 6, &records, 9, false);
AddKeyInfo( &records.value.abis.value, 6, &records.key, 6, false);
AddKeyInfo( &records.value.abis.value, 6, &records.value.abis, 9, false);
AddKeyInfo( &records.value.abis.value, 6, &records.value.abis.key, 5, false);
AddKeyInfo( &ADDR_INTERFACE_ID, 6, &ADDR_INTERFACE_ID, 9, false);
AddKeyInfo( &PUBKEY_INTERFACE_ID, 6, &PUBKEY_INTERFACE_ID, 9, false);
AddKeyInfo( &ens, 7, &ens, 9, false);
AddKeyInfo( &records.value.pubkey.y, 6, &records, 9, false);
AddKeyInfo( &records.value.pubkey.y, 6, &records.key, 6, false);
AddKeyInfo( &records.value.pubkey.y, 6, &records.value.pubkey, 9, false);
AddKeyInfo( &records.value.pubkey.y, 6, &records.value.pubkey.y, 9, false);
AddKeyInfo( &TEXT_INTERFACE_ID, 6, &TEXT_INTERFACE_ID, 9, false);
AddKeyInfo( &INTERFACE_META_ID, 6, &INTERFACE_META_ID, 9, false);
AddKeyInfo( &ABI_INTERFACE_ID, 6, &ABI_INTERFACE_ID, 9, false);
AddKeyInfo( &records.value.text.value, 6, &records, 9, false);
AddKeyInfo( &records.value.text.value, 6, &records.key, 6, false);
AddKeyInfo( &records.value.text.value, 6, &records.value.text, 9, false);
AddKeyInfo( &records.value.text.value, 6, &records.value.text.key, 6, false);
AddKeyInfo( &records.value.name, 6, &records, 9, false);
AddKeyInfo( &records.value.name, 6, &records.key, 6, false);
AddKeyInfo( &records.value.name, 6, &records.value.name, 9, false);
AddKeyInfo( &NAME_INTERFACE_ID, 6, &NAME_INTERFACE_ID, 9, false);
}
constructor PublicResolver(address ensAddr)
{
keyvbvhiap7();
InitializeVariables();
    ens = ensAddr;
}

void onlyOwner(string node)
{
    CallParams params = {ens, U256(0), 100000};
    address addr = owner(params, node);
    address sender = GetSender();
    if (Equal(addr, sender) != true)
    {
        Revert("Not Sender");
    }
}

   
                                                                                        
                                                           
                                                                   
   
UNMUTABLE
bool supportsInterface(string interfaceID)
{
keyvbvhiap7();
    return Equal(interfaceID, ADDR_INTERFACE_ID) ||
           Equal(interfaceID, CONTENT_INTERFACE_ID) ||
           Equal(interfaceID, NAME_INTERFACE_ID) ||
           Equal(interfaceID, ABI_INTERFACE_ID) ||
           Equal(interfaceID, PUBKEY_INTERFACE_ID) ||
           Equal(interfaceID, TEXT_INTERFACE_ID) ||
           Equal(interfaceID, INTERFACE_META_ID);
}

   
                                                   
                                     
                                  
   
MUTABLE
address addr(string node)
{
keyvbvhiap7();
    records.key = node;
    Record record = records.value;
    return record.addr;
}

   
                                                
                                                                    
                                  
                                  
   
MUTABLE
void setAddr(string node, address addr)
{
keyvbvhiap7();
    onlyOwner(node);
    records.key = node;
    records.value.addr = addr;
    AddrChanged(node, addr);
}

   
                                                        
                                                                           
                                                   
                                     
                                       
   
UNMUTABLE
string content(string node)
{
keyvbvhiap7();
    records.key = node;
    Record record = records.value;
    return record.content;
}

   
                                                     
                                                                    
                                                                           
                                                   
                                  
                                      
   
MUTABLE
void setContent(string node, string hash)
{
keyvbvhiap7();
    onlyOwner(node);
    records.key = node;
    records.value.content = hash;
    ContentChanged(node, hash);
}

   
                                                                     
                     
                                     
                               
   
UNMUTABLE
string name(string node)
{
keyvbvhiap7();
    records.key = node;
    Record record = records.value;
    return record.name;
}

   
                                                                  
                                                                    
                                  
                               
   
MUTABLE
void setName(string node, string name)
{
keyvbvhiap7();
    onlyOwner(node);
    records.key = node;
    records.value.name = name;
    NameChanged(node, name);
}

   
                                               
                     
                                    
                                                                              
                                                           
                            
   
UNMUTABLE
string ABIRecord(string node, uint256 contentTypes)
{
keyvbvhiap7();
    records.key = node;
    Record record = records.value;
    for (uint256 contentType = U256(1); U256_Cmp(contentType, contentTypes) != 1; contentType = U256_Shl(contentType, 1))
    {
        record.abis.key = contentType;
        string recordabis = record.abis.value;
        if (U256_Cmp(U256_And(contentType, contentTypes), U256(0)) != 0 && Equal(recordabis, "") == false)
        {
            return recordabis;
        }
    }
    return "";
}

UNMUTABLE
uint256 ABIContentType(string node, uint256 contentTypes)
{
keyvbvhiap7();
    records.key = node;
    Record record = records.value;
    for (uint256 contentType = U256(1); U256_Cmp(contentType, contentTypes) != 1; contentType = U256_Shl(contentType, 1))
    {
        record.abis.key = contentType;
        string recordabis = record.abis.value;
        if (U256_Cmp(U256_And(contentType, contentTypes), U256(0)) != 0 && Equal(recordabis, "") == false)
        {
            return contentType;
        }
    }
    return U256(0);
}

   
                                            
                                                                           
                    
                                  
                                                 
                            
   
MUTABLE
void setABI(string node, uint256 contentType, string data)
{
keyvbvhiap7();
    onlyOwner(node);
                                        
    if (U256_Cmp(U256_Add(U256_Sub(contentType, 1), contentType), 0) != 0)
    {
        Revert("");
    }
    records.key = node;
    records.value.abis.key = contentType;
    records.value.abis.value = data;
    ABIChanged(node, contentType);
}

   
                                                                
                      
                                    
                                                                              
   
UNMUTABLE
string pubkeyX(string node)
{
keyvbvhiap7();
    records.key = node;
    return records.value.pubkey.x;
}

UNMUTABLE
string pubkeyY(string node)
{
keyvbvhiap7();
    records.key = node;
    return records.value.pubkey.y;
}

   
                                                             
                                    
                                                                   
                                                                   
   
MUTABLE
void setPubkey(string node, string x, string y)
{
keyvbvhiap7();
    onlyOwner(node);
    records.key = node;
    records.value.pubkey.x = x;
    records.value.pubkey.y = y;
    PubkeyChanged(node, x, y);
}

   
                                                             
                                     
                                         
                                    
   
UNMUTABLE
string text(string node, string key)
{
keyvbvhiap7();
    records.key = node;
    records.value.text.key = key;
    return records.value.text.value;
}

   
                                                          
                                                                    
                                  
                             
                                           
   
MUTABLE
void setText(string node, string key, string value)
{
keyvbvhiap7();
    onlyOwner(node);
    records.key = node;
    records.value.text.key = key;
    records.value.text.value = value;
    TextChanged(node, key, key);
}

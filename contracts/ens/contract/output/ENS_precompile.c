#include "./vntlib.h"

typedef struct
{
    address owner;
    address resolver;
    uint64 ttl;
} Record;

KEY mapping(string, Record) records;

                                                                    
EVENT NewOwner(indexed string node, indexed string label, address owner);

                                                                        
EVENT Transfer(indexed string node, address owner);

                                                  
EVENT NewResolver(indexed string node, address resolver);

                                           
EVENT NewTTL(indexed string node, uint64 ttl);

   
                                  
   

void keyddslbvum(){
AddKeyInfo( &records.value.owner, 7, &records, 9, false);
AddKeyInfo( &records.value.owner, 7, &records.key, 6, false);
AddKeyInfo( &records.value.owner, 7, &records.value.owner, 9, false);
AddKeyInfo( &records.value.resolver, 7, &records, 9, false);
AddKeyInfo( &records.value.resolver, 7, &records.key, 6, false);
AddKeyInfo( &records.value.resolver, 7, &records.value.resolver, 9, false);
AddKeyInfo( &records.value.ttl, 4, &records, 9, false);
AddKeyInfo( &records.value.ttl, 4, &records.key, 6, false);
AddKeyInfo( &records.value.ttl, 4, &records.value.ttl, 9, false);
}
constructor ENS()
{
keyddslbvum();
InitializeVariables();
    records.key = "";
    records.value.owner = GetSender();
}

   
                                                    
   
                                                            
                                  
    

void onlyOwner(string node)
{
    records.key = node;
    address owner = records.value.owner;
    address sender = GetSender();
    if (Equal(owner, sender) == false)
    {
        Revert("need owner");
    }
}

UNMUTABLE
address owner(string node)
{
keyddslbvum();
    records.key = node;
    return records.value.owner;
}

   
                                                              
   
                                                               
                                     
    

UNMUTABLE
address resolver(string node)
{
keyddslbvum();
    records.key = node;
    return records.value.resolver;
}

   
                                                                 
   
                                                         
                                
    
UNMUTABLE
uint64 ttl(string node)
{
keyddslbvum();
    records.key = node;
    return records.value.ttl;
}

   
                                                                                    
                     
                                                 
                                             
   
                                                                    
                             
                                   
    

MUTABLE
void setOwner(string node, address owner)
{
keyddslbvum();
    onlyOwner(node);
    Transfer(node, owner);
    records.key = node;
    records.value.owner = owner;
}

   
                                                                                   
                                          
                               
                                                             
                                             
   
                                                                                          
                                       
                                    
                                      
    
MUTABLE
void setSubnodeOwner(string node, string label, address owner)
{
keyddslbvum();
    onlyOwner(node);
    string subnode = SHA3(Concat(node, label));
    NewOwner(node, label, owner);
    records.key = subnode;
    records.value.owner = owner;
}

   
                                                    
                                  
                                               
   
                                                                          
                                   
                                         
    

MUTABLE
void setResolver(string node, address resolver)
{
keyddslbvum();
    onlyOwner(node);
    NewResolver(node, resolver);
    records.key = node;
    records.value.resolver = resolver;
}

   
                                       
                                  
                                 
   
                                                               
                         
                               
    

MUTABLE
void setTTL(string node, uint64 ttl)
{
keyddslbvum();
    onlyOwner(node);
    records.key = node;
    records.value.ttl = ttl;
}
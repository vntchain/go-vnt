#include "./vntlib.h"

                           

                         

                                              
                                                  
                                  
                                                           
                                                

                                  
                                         

                                                
                                    

                              
          
                                                 
                                                 
                                             
                                             
                                             
                                                                                                                                         
                                                                                                             
                                         
                                                                                        
                                               
                                                        
                                                                        
                                                                  
                                                                                        
                                                                                 
                                                     
                                      
                                                             
                                          
                                          
                   
                                                                          
                                              
                                
                                         
                                         
            
        
    

KEY mapping(address, uint256) sent;
KEY address owner;

EVENT Overdraft(address deadbeat);


void keyrxv22r0m(){
AddKeyInfo( &sent.value, 5, &sent, 9, false);
AddKeyInfo( &sent.value, 5, &sent.key, 7, false);
AddKeyInfo( &owner, 7, &owner, 9, false);
}
constructor chequebook()
{
keyrxv22r0m();
InitializeVariables();
    owner = GetSender();
}

MUTABLE
void cash(address beneficiary, uint256 amount, string sig_v, string sig_r, string sig_s)
{
keyrxv22r0m();
    sent.key = beneficiary;
    Require(U256_Cmp(amount, sent.value) == 1, "amount is too small");
    string hash = SHA3(Concat(Concat(GetContractAddress(), beneficiary), U256ToString(amount)));
    address recover = Ecrecover(hash, sig_v, sig_r, sig_s);
    PrintAddress("owner:", owner);
    PrintAddress("recover:", recover);
    Require(Equal(owner, recover), "is not owner");
    uint256 diff = U256SafeSub(amount, sent.value);
    uint256 balance = GetBalanceFromAddress(GetContractAddress());
    if (U256_Cmp(diff, balance) != 1)
    {
        sent.key = beneficiary;
        sent.value = amount;
        SendFromContract(beneficiary, diff);
    }
    else
    {
        Overdraft(owner);
        SendFromContract(beneficiary, balance);
    }
}

MUTABLE
void kill()
{
keyrxv22r0m();
}

UNMUTABLE
uint256 GetSent(address beneficiary)
{
keyrxv22r0m();
    sent.key = beneficiary;
    return sent.value;
}

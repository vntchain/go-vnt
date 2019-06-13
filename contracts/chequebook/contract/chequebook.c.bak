// +build none
#include "./vntlib.h"

// pragma solidity ^0.4.18;

// import "./mortal.sol";

// /// @title Chequebook for VNT micropayments
// /// @author Daniel A. Nagy <daniel@vntchain.io>
// contract chequebook is mortal {
//     // Cumulative paid amount in wei to each beneficiary
//     mapping (address => uint256) public sent;

//     /// @notice Overdraft event
//     event Overdraft(address deadbeat);

//     // Allow sending ether to the chequebook.
//     function() public payable { }

//     /// @notice Cash cheque
//     ///
//     /// @param beneficiary beneficiary address
//     /// @param amount cumulative amount in wei
//     /// @param sig_v signature parameter v
//     /// @param sig_r signature parameter r
//     /// @param sig_s signature parameter s
//     /// The digital signature is calculated on the concatenated triplet of contract address, beneficiary address and cumulative amount
//     function cash(address beneficiary, uint256 amount, uint8 sig_v, bytes32 sig_r, bytes32 sig_s) public {
//         // Check if the cheque is old.
//         // Only cheques that are more recent than the last cashed one are considered.
//         require(amount > sent[beneficiary]);
//         // Check the digital signature of the cheque.
//         bytes32 hash = keccak256(address(this), beneficiary, amount);
//         require(owner == ecrecover(hash, sig_v, sig_r, sig_s));
//         // Attempt sending the difference between the cumulative amount on the cheque
//         // and the cumulative amount on the last cashed cheque to beneficiary.
//         uint256 diff = amount - sent[beneficiary];
//         if (diff <= this.balance) {
//             // update the cumulative amount before sending
//             sent[beneficiary] = amount;
//             beneficiary.transfer(diff);
//         } else {
//             // Upon failure, punish owner for writing a bounced cheque.
//             // owner.sendToDebtorsPrison();
//             Overdraft(owner);
//             // Compensate beneficiary.
//             selfdestruct(beneficiary);
//         }
//     }
// }

KEY mapping(address, uint256) sent;
KEY address owner;

EVENT Overdraft(address deadbeat);

constructor $chequebook()
{
    owner = GetSender();
}

MUTABLE
void cash(address beneficiary, uint256 amount, string sig_v, string sig_r, string sig_s)
{
    sent.key = beneficiary;
    Require(U256_Cmp(amount, sent.value) == 1, "amount is too small");
    string hash = SHA3(Concat(Concat(GetContractAddress(), beneficiary), U256ToString(amount)));
    address recover = Ecrecover(hash, sig_v, sig_r, sig_s);
    PrintAddress("owner:", owner);
    PrintAddress("recover:", recover);
    Require(Equal(owner, recover) == true, "is not owner");
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
}

UNMUTABLE
uint256 GetSent(address beneficiary)
{
    sent.key = beneficiary;
    return sent.value;
}

$_()
{
}

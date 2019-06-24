#include "../../../vntlib/vntlib.h"

KEY mapping(address, uint256) sent;
KEY address owner;

EVENT Overdraft(address deadbeat);

constructor $chequebook() { owner = GetSender(); }

MUTABLE
void cash(address beneficiary, uint256 amount, string sig_v, string sig_r,
          string sig_s) {
  sent.key = beneficiary;
  Require(U256_Cmp(amount, sent.value) == 1, "amount is too small");
  string hash = SHA3(
      Concat(Concat(GetContractAddress(), beneficiary), U256ToString(amount)));
  address recover = Ecrecover(hash, sig_v, sig_r, sig_s);
  PrintAddress("owner:", owner);
  PrintAddress("recover:", recover);
  Require(Equal(owner, recover) == true, "is not owner");
  uint256 diff = U256SafeSub(amount, sent.value);
  uint256 balance = GetBalanceFromAddress(GetContractAddress());
  if (U256_Cmp(diff, balance) != 1) {
    sent.key = beneficiary;
    sent.value = amount;
    SendFromContract(beneficiary, diff);
  } else {
    Overdraft(owner);
    SendFromContract(beneficiary, balance);
  }
}

MUTABLE
void kill() {}

UNMUTABLE
uint256 GetSent(address beneficiary) {
  sent.key = beneficiary;
  return sent.value;
}

$_() {}

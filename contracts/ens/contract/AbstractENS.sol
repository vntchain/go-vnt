pragma solidity ^0.4.0;

contract AbstractENS {
    function owner(string node) constant returns(address);
    function resolver(string node) constant returns(address);
    function ttl(string node) constant returns(uint64);
    function setOwner(string node, address owner);
    function setSubnodeOwner(string node, string label, address owner);
    function setResolver(string node, address resolver);
    function setTTL(string node, uint64 ttl);

    // Logged when the owner of a node assigns a new owner to a subnode.
    event NewOwner(string indexed node, string indexed label, address owner);

    // Logged when the owner of a node transfers ownership to a new account.
    event Transfer(string indexed node, address owner);

    // Logged when the resolver for a node changes.
    event NewResolver(string indexed node, address resolver);

    // Logged when the TTL of a node changes
    event NewTTL(string indexed node, uint64 ttl);
}

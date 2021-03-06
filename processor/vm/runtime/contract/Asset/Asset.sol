pragma solidity ^0.4.24;

contract Asset {
    uint256 public totalSupply;

    constructor() public {
        totalSupply = 10;
    }

    function reg(string desc) public {
        issueasset(desc);
    }
    function add(address assetId, uint256 value) public {
        addasset(assetId,value);
    }
    function transAsset(address to, address assetId, uint256 value) public payable {
        to.transferex(assetId, value);
    }
    function setname(address newOwner, address assetId) public {
        setassetowner(assetId, newOwner);
    }
    function getbalance(address to, address assetId) public returns(uint) {
        return to.balanceex(assetId);
    }
    function getAssetAmount(uint256 assetId, uint256 t) public returns (uint256){
        return assetamount(assetId,t);
    }
    function getSnapshotTime(uint256 i,uint256 t) public returns (uint256){
        return snapshottime(i,t);
    }
    function getSnapBalance(address to,uint256 assetId,uint256 t) public returns (uint256){
        return to.snapbalance(assetId,t);
    }
}

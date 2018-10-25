// Copyright (c) 2009-2010 Satoshi Nakamoto
// Copyright (c) 2009-2017 The Bitcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#ifndef BITCOIN_PRIMITIVES_BLOCK_H
#define BITCOIN_PRIMITIVES_BLOCK_H

#include <primitives/transaction.h>
#include <serialize.h>
#include <uint256.h>

/** Nodes collect new transactions into a block, hash them into a hash tree,
 * and scan through nonce values to make the block's hash satisfy proof-of-work
 * requirements.  When they solve the proof-of-work, they broadcast the block
 * to everyone and the block is added to the block chain.  The first transaction
 * in the block is a special one that creates a new coin owned by the creator
 * of the block.
 */
class CBlockHeader
{
public:
    // header
    int32_t				nVersion;
	int32_t				nHeight;
    uint256				hashPrevBlock;
    uint256				hashMerkleRoot;
    uint32_t			nTime;
    uint32_t			nBits;
    uint32_t			nNonce;

	mutable uint256     hashSeed_btcv;
	mutable uint256		hashLock_btcv;
	mutable uint32_t	nNonceLock_btcv;
	mutable uint256		hashLockSeed_btcv;
	mutable uint256		hashBlockSeed_btcv;
	mutable uint256		hashBlock_btcv;
	mutable uint32_t	nColNum_btcv;

	CBlockHeader()
    {
        SetNull();
    }

    ADD_SERIALIZE_METHODS;

    template <typename Stream, typename Operation>
    inline void SerializationOp(Stream& s, Operation ser_action) {
        READWRITE(this->nVersion);
		READWRITE(this->nHeight);
        READWRITE(hashPrevBlock);
        READWRITE(hashMerkleRoot);
        READWRITE(nTime);
        READWRITE(nBits);
        READWRITE(nNonce);
		READWRITE(hashSeed_btcv);
		READWRITE(hashLock_btcv);
		READWRITE(nNonceLock_btcv);
		READWRITE(hashLockSeed_btcv);
		READWRITE(hashBlockSeed_btcv);
		READWRITE(hashBlock_btcv);
    }

    void SetNull()
    {
        nVersion = 0;
        nHeight = 0;
		hashPrevBlock.SetNull();
        hashMerkleRoot.SetNull();
        nTime = 0;
        nBits = 0;
        nNonce = 0;

		hashSeed_btcv.SetNull();
		hashLock_btcv.SetNull();
		hashLockSeed_btcv.SetNull();
		hashBlockSeed_btcv.SetNull();
		hashBlock_btcv.SetNull();
		nNonceLock_btcv = 0;
    }

    bool IsNull() const
    {
        return (nBits == 0);
    }

    uint256 GetHash( bool abCheck = true ) const;
	void	FillingSeedByNonce( uint32_t aui32Nonce, 
								unsigned char* apHashIn, 
								unsigned char* apHashOut ) const;

	bool	CheckSame(			unsigned char* apLeft, 
								unsigned char* apRight, 
								int aiMove )  const;


	//////////////// in /////////////////////
	void InitBlockLock( void ) const;


    int64_t GetBlockTime() const
    {
        return (int64_t)nTime;
    }
};


class CBlock : public CBlockHeader
{
public:
    // network and disk
    std::vector<CTransactionRef> vtx;

    // memory only
    mutable bool fChecked;

    CBlock()
    {
        SetNull();
    }

    CBlock(const CBlockHeader &header)
    {
        SetNull();
        *((CBlockHeader*)this) = header;
    }

    ADD_SERIALIZE_METHODS;

    template <typename Stream, typename Operation>
    inline void SerializationOp(Stream& s, Operation ser_action) {
        READWRITE(*(CBlockHeader*)this);
        READWRITE(vtx);
    }

    void SetNull()
    {
        CBlockHeader::SetNull();
        vtx.clear();
        fChecked = false;
    }

    CBlockHeader GetBlockHeader() const
    {
        CBlockHeader block;
        block.nVersion             = nVersion;
		block.nHeight		       = nHeight;
        block.hashPrevBlock        = hashPrevBlock;
        block.hashMerkleRoot       = hashMerkleRoot;
        block.nTime                = nTime;
        block.nBits                = nBits;
        block.nNonce			   = nNonce;
		block.hashSeed_btcv         = hashSeed_btcv;
		block.hashLock_btcv         = hashLock_btcv;
		block.hashLockSeed_btcv     = hashLockSeed_btcv;
		block.hashBlockSeed_btcv    = hashBlockSeed_btcv;
		block.hashBlock_btcv        = hashBlock_btcv;
		block.nNonceLock_btcv       = nNonceLock_btcv;
        return block;
    }

    std::string ToString() const;
};

/** Describes a place in the block chain to another node such that if the
 * other node doesn't have the same branch, it can find a recent common trunk.
 * The further back it is, the further before the fork it may be.
 */
struct CBlockLocator
{
    std::vector<uint256> vHave;

    CBlockLocator() {}

    explicit CBlockLocator(const std::vector<uint256>& vHaveIn) : vHave(vHaveIn) {}

    ADD_SERIALIZE_METHODS;

    template <typename Stream, typename Operation>
    inline void SerializationOp(Stream& s, Operation ser_action) {
        int nVersion = s.GetVersion();
        if (!(s.GetType() & SER_GETHASH))
            READWRITE(nVersion);
        READWRITE(vHave);
    }

    void SetNull()
    {
        vHave.clear();
    }

    bool IsNull() const
    {
        return vHave.empty();
    }
};

#endif // BITCOIN_PRIMITIVES_BLOCK_H

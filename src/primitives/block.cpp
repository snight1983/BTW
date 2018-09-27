// Copyright (c) 2009-2010 Satoshi Nakamoto
// Copyright (c) 2009-2017 The Bitcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include <primitives/block.h>
#include <hash.h>
#include <tinyformat.h>
#include <utilstrencodings.h>
#include <crypto/common.h>

#define							DEF_32_LEN							( 32 )
#define							DEF_COLUMN_CNT						( 8192 )
#define							DEF_OVERALL_SIZE					( 262144 )

bool CBlockHeader::CheckSame( unsigned int* apLeft, unsigned int* apRight, int aiMove ) const {
	return ( ( (*apLeft) << aiMove) == ((*apRight) << aiMove) ) ? true:false;
}

void CBlockHeader::FillingSeedByNonce( uint32_t aui32Nonce, unsigned char* apHashIn, unsigned char* apHashOut) const
{

}

bool CBlockHeader::GetHashLock() const{

}

uint256 CBlockHeader::GetHash( bool abCheck ) const{

	return hashBlock_btw;
}

std::string CBlock::ToString() const
{
    std::stringstream s;
    s << strprintf("CBlock(hash=%s, ver=0x%08x, hashPrevBlock=%s, hashMerkleRoot=%s, nTime=%u, nBits=%08x, nNonce=%u, vtx=%u)\n",
        GetHash().ToString(),
        nVersion,
        hashPrevBlock.ToString(),
        hashMerkleRoot.ToString(),
        nTime, nBits, nNonce,
        vtx.size());
    for (const auto& tx : vtx) {
        s << "  " << tx->ToString() << "\n";
    }
    return s.str();
}

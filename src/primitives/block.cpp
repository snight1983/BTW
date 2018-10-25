// Copyright (c) 2009-2010 Satoshi Nakamoto
// Copyright (c) 2009-2017 The Bitcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include <primitives/block.h>
#include <hash.h>
#include <tinyformat.h>
#include <utilstrencodings.h>
#include <crypto/common.h>

#define	DEF_32_LEN					( 32 )
#define	DEF_COLUMN_CNT				( 8192 )
#define	DEF_OVERALL_SIZE			( 262144 )

bool CBlockHeader::CheckSame( unsigned char* apLeft, unsigned char* apRight, int aiMove ) const {
	for ( int i =0; i< aiMove; ++i ){
		if ( * (apLeft+i) != *(apRight+i)){
			return false;
		}
	}
	return true;
}


void CBlockHeader::InitBlockLock(void) const{
}

void CBlockHeader::FillingSeedByNonce( uint32_t aui32Nonce, unsigned char* apHashIn, unsigned char* apHashOut) const
{
	unsigned char lszOverall[DEF_OVERALL_SIZE+1]			=	{0};
	unsigned char lszHashNum[ 43]							=	{0};
	memcpy( lszHashNum, apHashIn, DEF_32_LEN );
	char lsznNonce[11]										=	{0};
	sprintf( lsznNonce,"%u",  aui32Nonce );
	int liNonceLen											=	strlen(lsznNonce);
	memcpy( lszHashNum + DEF_32_LEN, lsznNonce, liNonceLen );
	unsigned char lszHashOut[DEF_32_LEN+1]					=	{0};
	CryptoVIP(lszHashOut, lszHashNum, DEF_32_LEN + liNonceLen);
	for( int i = 0; i < DEF_COLUMN_CNT; ++i ){
		memset( lszHashNum, 0, 43);
		memcpy( lszHashNum, lszHashOut, DEF_32_LEN );
		memcpy( lszHashNum + DEF_32_LEN, lsznNonce, liNonceLen );
		CryptoVIP(lszHashOut, lszHashNum, DEF_32_LEN + liNonceLen );
		memcpy( lszOverall+i*DEF_32_LEN, lszHashOut, DEF_32_LEN );
	}
	CryptoVIP(lszHashOut, lszOverall, DEF_OVERALL_SIZE );
	unsigned int liAry[8] = {0};
	for( int j = 0; j < 8; ++j ) {
		memcpy(&liAry[j], lszHashOut+j*4, 4 );
	}
	for( int k = 0; k < 8; ++k ){
		unsigned int liValue;
		memcpy( &liValue, lszOverall + nColNum_btcv*DEF_32_LEN+k*4, 4 );
		liValue = liValue^liAry[k];
		memcpy( apHashOut + k*4, &liValue, 4);
	}
}

uint256 CBlockHeader::GetHash( bool abCheck ) const{
	nColNum_btcv = hashPrevBlock.GetUint64(0) % DEF_COLUMN_CNT;
	
	return hashBlock_btcv;
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

// Copyright (c) 2009-2010 Satoshi Nakamoto
// Copyright (c) 2009-2017 The Bitcoin Core developers
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

#include <primitives/block.h>

#include <hash.h>
#include <tinyformat.h>
#include <utilstrencodings.h>
#include <crypto/common.h>

#define							DEF_HASH_LEN						( 32 )
#define							DEF_COLUMN_CNT						( 8192 )
#define							DEF_OVERALL_SIZE					( 262144 )

void CBlockHeader::FillingSeedByNonce( uint32_t aui32Nonce, unsigned char* apHashIn, unsigned char* apHashOut) const
{
	unsigned char lszPHash[256]								=	{0};
	unsigned char lszOverall[DEF_OVERALL_SIZE]				=	{0};
	memcpy( lszPHash, apHashIn, DEF_HASH_LEN );
	char lsznNonce[256]										=	{0};
	sprintf( lsznNonce,"%u",  aui32Nonce );
	int liNonceLen											=	strlen(lsznNonce);
	memcpy( lszPHash + DEF_HASH_LEN, lsznNonce, DEF_HASH_LEN );

	unsigned char lszHashBuf[DEF_HASH_LEN]					=	{0};

	CSHA256 loSha256;
	loSha256.Reset();
	loSha256.Write( lszPHash,  DEF_HASH_LEN + liNonceLen );
	loSha256.Finalize( lszHashBuf );

	unsigned char lszBuf[ 500 ]		=	{0};
	for( int i = 0; i < DEF_COLUMN_CNT; ++i ){
		loSha256.Reset();
		memcpy( lszBuf, lszHashBuf, DEF_HASH_LEN );
		memcpy( lszBuf + DEF_HASH_LEN, lsznNonce, liNonceLen );
		loSha256.Write( lszBuf,  DEF_HASH_LEN + liNonceLen );
		loSha256.Finalize( lszHashBuf );
		memcpy( lszOverall+i*DEF_HASH_LEN, lszHashBuf, DEF_HASH_LEN );
	}
	loSha256.Reset();
	loSha256.Write( lszOverall,  DEF_OVERALL_SIZE );
	loSha256.Finalize( lszHashBuf );

	unsigned int liAry[8] = {0};
	for( int j = 0; j < 8; ++j ){
		memcpy(&liAry[j], lszHashBuf+j*4, 4 );
	}
	for( int k = 0; k < 8; ++k ){
		unsigned int liValue;
		memcpy( &liValue, lszOverall + nColNum_btw*DEF_HASH_LEN+k*4, 4 );
		liValue = liValue^liAry[k];
		memcpy( apHashOut + k*4, &liValue, 4);
	}
}


bool CBlockHeader::CheckSame( unsigned int* apLeft, unsigned int* apRight, int aiMove ) const {
	return ( ( (*apLeft) << aiMove) == ((*apRight) << aiMove) ) ? true:false;
}

bool CBlockHeader::GetHashLock() const{

}

uint256 CBlockHeader::GetHash( bool abCheck ) const
{
	nColNum_btw = hashPrevBlock.GetUint64(0) % DEF_COLUMN_CNT;
	/**/
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

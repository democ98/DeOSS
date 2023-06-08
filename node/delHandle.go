/*
	Copyright (C) CESS. All rights reserved.
	Copyright (C) Cumulus Encrypted Storage System. All rights reserved.

	SPDX-License-Identifier: Apache-2.0
*/

package node

import (
	"fmt"
	"net/http"
	"unsafe"

	"github.com/CESSProject/sdk-go/core/pattern"
	sutils "github.com/CESSProject/sdk-go/core/utils"
	"github.com/gin-gonic/gin"
)

type DelList struct {
	FileId []string `json:"file_id"`
}

// delHandle is used to delete buckets or files
func (n *Node) delHandle(c *gin.Context) {
	var (
		err      error
		txHash   string
		clientIp string
		pkey     []byte
		respMsg  = &RespMsg{}
	)

	clientIp = c.ClientIP()
	n.Upfile("info", fmt.Sprintf("[%v] %v", clientIp, INFO_DelRequest))

	// verify token
	_, pkey, err = n.VerifyToken(c, respMsg)
	if err != nil {
		n.Upfile("err", fmt.Sprintf("[%v] %v", clientIp, err))
		c.JSON(respMsg.Code, respMsg.Err)
		return
	}

	deleteName := c.Param(PUT_ParameterName)
	if len(deleteName) == int(unsafe.Sizeof(pattern.FileHash{})) {
		txHash, _, err = n.DeleteFile(pkey, []string{deleteName})
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, txHash)
	} else if sutils.CheckBucketName(deleteName) {
		txHash, err = n.DeleteBucket(pkey, deleteName)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}
		c.JSON(200, txHash)
	} else {
		deleteNames := c.PostFormArray("delete_list")
		if err != nil {
			c.JSON(400, "InvalidBody.DeleteList")
			return
		}

		txHash, failList, err := n.DeleteFile(pkey, deleteNames)
		if err != nil {
			c.JSON(400, err.Error())
			return
		}

		c.JSON(http.StatusOK, struct {
			Block_hash  string
			Failed_list []pattern.FileHash
		}{Block_hash: txHash, Failed_list: failList})
		return
	}
}

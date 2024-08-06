/**
 *
 * Copyright © 2015--2018 . All rights reserved.
 *
 * File: address_test
 * Date: 2018-09-04
 *
 */

package ligo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNewPrivate(t *testing.T) {
	priv, pub, addr, err := GetNewPrivate()
	assert.Nil(t, err)
	fmt.Println("private:", priv)
	fmt.Println("pubkey:", pub)
	fmt.Println("address:", addr)
}
func TestGetAddressBytes(t *testing.T) {

}

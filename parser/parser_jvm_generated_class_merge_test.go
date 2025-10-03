package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeJVMGeneratedClasses(t *testing.T) {
	testcases := []struct {
		src         string
		expectedRes string
	}{
		{
			"org/example/rideshare/EnclosingClass$$Lambda$4/1283928880",
			"org/example/rideshare/EnclosingClass$$Lambda_",
		},
		{
			"org/example/rideshare/EnclosingClass$$Lambda$8/0x0000000800c01220",
			"org/example/rideshare/EnclosingClass$$Lambda_",
		},
		{
			"Fib$$Lambda.0x00007ffa600c4da0",
			"Fib$$Lambda_",
		},
		{
			"java/util/concurrent/Executors$RunnableAdapter",
			"java/util/concurrent/Executors$RunnableAdapter",
		},
		{
			"jdk/internal/reflect/GeneratedMethodAccessor31",
			"jdk/internal/reflect/GeneratedMethodAccessor_",
		},
		{
			"foo/bar/Baz$$EnhancerBySpringCGLIB$$1234567890",
			"foo/bar/Baz$$EnhancerBySpringCGLIB$$_",
		},
	}
	for _, testcase := range testcases {
		res := mergeJVMGeneratedClasses(testcase.src)
		assert.Equal(t, testcase.expectedRes, res)
	}
}

func TestMergeSharedLibs(t *testing.T) {

	testcases := []struct {
		src         string
		expectedRes string
	}{
		{
			"libasyncProfiler-linux-arm64-17b9a1d8156277a98ccc871afa9a8f69215f92.so",
			"libasyncProfiler-_.so",
		},
		{
			"libasyncProfiler-linux-musl-x64-17b9a1d8156277a98ccc871afa9a8f69215f92.so",
			"libasyncProfiler-_.so",
		},
		{
			"libasyncProfiler-linux-x64-17b9a1d8156277a98ccc871afa9a8f69215f92.so",
			"libasyncProfiler-_.so",
		},
		{
			"libasyncProfiler-macos-17b9a1d8156277a98ccc871afa9a8f69215f92.so",
			"libasyncProfiler-_.so",
		},
		{
			"libamazonCorrettoCryptoProvider109b39cf33c563eb.so",
			"libamazonCorrettoCryptoProvider_.so",
		},
		{
			"amazonCorrettoCryptoProviderNativeLibraries.7382c2f79097f415/libcrypto.so",
			"libamazonCorrettoCryptoProvider_.so",
		},
		{
			"amazonCorrettoCryptoProviderNativeLibraries.24e42b0d5ecf5f50/libamazonCorrettoCryptoProvider.so",
			"libamazonCorrettoCryptoProvider_.so",
		},
		{
			"libzstd-jni-1.5.1-16931311898282279136.so",
			"libzstd-jni-_.so",
		},
	}
	for _, testcase := range testcases {
		res := mergeJVMGeneratedClasses(testcase.src)
		assert.Equal(t, testcase.expectedRes, res)

		res = mergeJVMGeneratedClasses(testcase.src + " (deleted)")
		assert.Equal(t, testcase.expectedRes, res)

		res = mergeJVMGeneratedClasses("/tmp/" + testcase.src + " (deleted)")
		assert.Equal(t, testcase.expectedRes, res)

		res = mergeJVMGeneratedClasses("./tmp/" + testcase.src + " (deleted)")
		assert.Equal(t, testcase.expectedRes, res)
	}
}

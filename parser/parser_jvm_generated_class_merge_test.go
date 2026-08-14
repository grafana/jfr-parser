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
			"io/opentelemetry/context/Context$$Lambda$lambda$wrap$2$2630016632.0x00000000000100dd",
			"io/opentelemetry/context/Context$$Lambda$lambda$wrap$2_",
		},
		{
			"org/example/rideshare/EnclosingClass$$Lambda$lambda$doWork$1$1283928880",
			"org/example/rideshare/EnclosingClass$$Lambda$lambda$doWork$1_",
		},
		{
			"org/example/rideshare/EnclosingClass$$Lambda",
			"org/example/rideshare/EnclosingClass$$Lambda",
		},
		{
			"org/example/rideshare/EnclosingClass$$LambdaHelper",
			"org/example/rideshare/EnclosingClass$$LambdaHelper",
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

// A marker exists only to skip regex work, so it has to hold for everything its
// regex can match. A marker that does not follow from its regex disables the
// rule outright, with no other symptom.
func TestReplacementMarkersFollowFromRegex(t *testing.T) {
	samples := []string{
		"jdk/internal/reflect/GeneratedMethodAccessor31",
		"org/example/rideshare/EnclosingClass$$Lambda$8/0x0000000800c01220",
		"/tmp/libzstd-jni-1.5.1-16931311898282279136.so",
		"/tmp/libamazonCorrettoCryptoProvider109b39cf33c563eb.so",
		"/tmp/libasyncProfiler-macos-17b9a1d8156277a98ccc871afa9a8f69215f92.so",
		"foo/bar/Baz$$EnhancerBySpringCGLIB$$1234567890",
	}
	reached := make([]bool, len(replacements))
	for _, name := range samples {
		for i, r := range replacements {
			if !r.regex.MatchString(name) {
				continue
			}
			reached[i] = true
			assert.Equal(t, r.regex.ReplaceAllString(name, r.replaceWith), mergeJVMGeneratedClasses(name),
				"%q matches the regex of rule %q but merging did not apply it, so its marker skipped the rule",
				name, r.marker)
		}
	}
	for i, ok := range reached {
		assert.True(t, ok, "no sample reaches the regex of rule %q", replacements[i].marker)
	}
}

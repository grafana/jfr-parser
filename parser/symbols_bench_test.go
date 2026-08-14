package parser

import "testing"

// One name per replacement rule, in the order the rules are declared.
var benchAllRules = []string{
	"jdk/internal/reflect/GeneratedMethodAccessor31",
	"org/example/rideshare/EnclosingClass$$Lambda$8/0x0000000800c01220",
	"io/opentelemetry/context/Context$$Lambda$lambda$wrap$2$2630016632.0x00000000000100dd",
	"/tmp/libzstd-jni-1.5.1-16931311898282279136.so",
	"./tmp/amazonCorrettoCryptoProviderNativeLibraries.7382c2f79097f415/libcrypto.so (deleted)",
	"/tmp/libasyncProfiler-linux-arm64-17b9a1d8156277a98ccc871afa9a8f69215f92.so",
	"foo/bar/Baz$$EnhancerBySpringCGLIB$$1234567890",
}

// Ordinary frames, which are the overwhelming majority of a real symbol table.
// The last two carry a marker but no rewrite, so they pay for a regex that fails.
var benchUnmatched = []string{
	"java/lang/Thread.run",
	"com/bloom/catalog/CatalogService.search",
	"org/apache/tomcat/util/threads/ThreadPoolExecutor.runWorker",
	"org/springframework/web/filter/OncePerRequestFilter.doFilter",
	"java/util/regex/Pattern$CharPropertyGreedy.match",
	"java/util/concurrent/Executors$RunnableAdapter",
	"org/example/rideshare/EnclosingClass$$LambdaHelper",
	"/tmp/libzstd-jni-notaversion.so",
}

func BenchmarkMergeJVMGeneratedClasses(b *testing.B) {
	corpora := []struct {
		name  string
		names []string
	}{
		{"AllRules", benchAllRules},
		{"Unmatched", benchUnmatched},
	}
	for _, c := range corpora {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, name := range c.names {
					mergeJVMGeneratedClasses(name)
				}
			}
		})
	}
}

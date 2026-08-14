package parser

import (
	"regexp"
	"strings"

	"github.com/grafana/jfr-parser/parser/types"
)

type replacementRule struct {
	marker         string         // used to efficiently check if this rule applies
	markerIsPrefix bool           // used to specify if the marker is a prefix or can be anywhere in the string
	regex          *regexp.Regexp // the regex used for replacement matching
	replaceWith    string         // the actual replacement
}

// TODO
// ./tmp/snappy-1.1.8-6fb9393a-3093-4706-a7e4-837efe01d078-libsnappyjava.so
var replacements = []replacementRule{
	// jdk/internal/reflect/GeneratedMethodAccessor31
	{
		marker:         "jdk/internal/reflect/GeneratedMethodAccessor",
		markerIsPrefix: true,
		regex:          regexp.MustCompile(`^(jdk/internal/reflect/GeneratedMethodAccessor)(\d+)$`),
		replaceWith:    "${1}_",
	},
	// org/example/rideshare/OrderService$$Lambda$669.0x0000000800fd7318.run
	// Fib$$Lambda.0x00007ffa600c4da0.run
	// io/opentelemetry/context/Context$$Lambda$lambda$wrap$2$2630016632.0x00000000000100dd.call
	{
		marker:      "$$Lambda",
		regex:       regexp.MustCompile(`^(.+\$\$Lambda)(\$[a-zA-Z][\w$]*?)?((\$\d+)?[./](0x)?[\da-f]+|\$?\d+)$`),
		replaceWith: "${1}${2}_",
	},
	// libzstd-jni-1.5.1-16931311898282279136.so.Java_com_github_luben_zstd_ZstdInputStreamNoFinalizer_decompressStream
	{
		marker:      "libzstd-jni-",
		regex:       regexp.MustCompile(`^(\.?/tmp/)?(libzstd-jni-\d+\.\d+\.\d+-)(\d+)(\.so)( \(deleted\))?$`),
		replaceWith: "libzstd-jni-_.so",
	},
	// ./tmp/libamazonCorrettoCryptoProvider109b39cf33c563eb.so
	// ./tmp/amazonCorrettoCryptoProviderNativeLibraries.7382c2f79097f415/libcrypto.so (deleted)
	{
		marker: "amazonCorrettoCryptoProvider",
		regex: regexp.MustCompile(`^(\.?/tmp/)?(lib)?(amazonCorrettoCryptoProvider)(NativeLibraries\.)?([0-9a-f]{16})` +
			`(/libcrypto|/libamazonCorrettoCryptoProvider)?(\.so)( \(deleted\))?$`),
		replaceWith: "libamazonCorrettoCryptoProvider_.so",
	},
	// libasyncProfiler-linux-arm64-17b9a1d8156277a98ccc871afa9a8f69215f92.so
	{
		marker: "libasyncProfiler-",
		regex: regexp.MustCompile(
			`^(\.?/tmp/)?(libasyncProfiler)-(linux-arm64|linux-musl-x64|linux-x64|macos)-(17b9a1d8156277a98ccc871afa9a8f69215f92)(\.so)( \(deleted\))?$`),
		replaceWith: "libasyncProfiler-_.so",
	},
	// foo/bar/Baz$$EnhancerBySpringCGLIB$$1234567890
	{
		marker:      "$$EnhancerBySpringCGLIB$$",
		regex:       regexp.MustCompile(`^(.+\$\$EnhancerBySpringCGLIB\$\$)(.*)$`),
		replaceWith: "${1}_",
	},
}

func mergeJVMGeneratedClasses(frame string) string {
	// The marker check is written out here rather than behind a helper: as a
	// method it costs 102 inline nodes against a budget of 80, so every rule
	// would pay a real call on the path that exists to avoid work.
	for i := range replacements {
		r := &replacements[i]
		var hasMarker bool
		if r.markerIsPrefix {
			hasMarker = strings.HasPrefix(frame, r.marker)
		} else {
			hasMarker = strings.Contains(frame, r.marker)
		}
		if hasMarker {
			frame = r.regex.ReplaceAllString(frame, r.replaceWith)
		}
	}
	return frame
}

func ProcessSymbols(ref *types.SymbolList) {
	for i := range ref.Symbol { //todo regex replace inplace
		ref.Symbol[i].String = mergeJVMGeneratedClasses(ref.Symbol[i].String)
	}
}

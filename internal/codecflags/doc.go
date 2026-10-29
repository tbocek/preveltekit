// Package codecflags holds one constant per reflect kind supported by the
// JS codec (Decode/Encode). Each constant is true unless the build tag
// pk_no<kind> is set, in which case the codec drops support for that kind
// and the compiler removes the code for it.
//
// cmd/reflectscan analyzes an app and prints the tags for kinds it never
// decodes or encodes; build.sh passes them to tinygo.
package codecflags

package base62

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBase62(t *testing.T) {
	n := int64(1517153236107137025)
	s := "1O4zPvQMvmh"
	Convey("base62:Encode", t, func() {
		So(Encode(n), ShouldEqual, s)
	})

	Convey("base62:Decode", t, func() {
		So(Decode(s), ShouldEqual, n)
	})
}

func BenchmarkBase62Encode(b *testing.B) {
	n := int64(1517153236107137025)
	for i := 0; i < b.N; i++ {
		Encode(n)
	}
}

func BenchmarkBase62Decode(b *testing.B) {
	s := "1O4zPvQMvmh"
	for i := 0; i < b.N; i++ {
		Decode(s)
	}
}

func BenchmarkBase62EncodeParallel(b *testing.B) {
	n := int64(1517153236107137025)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Encode(n)
		}
	})
}

func BenchmarkBase62DecodeParallel(b *testing.B) {
	s := "1O4zPvQMvmh"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Decode(s)
		}
	})
}

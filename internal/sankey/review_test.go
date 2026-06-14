package sankey

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNegativeFlowRejected(t *testing.T) {
	Convey("Given a sankey flow with a negative value", t, func() {
		Convey("When parsing", func() {
			_, err := Parse("sankey-beta\nA,B,-5")
			Convey("Then it returns an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

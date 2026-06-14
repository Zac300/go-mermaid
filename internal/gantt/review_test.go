package gantt

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func findTask(d *Diagram, id string) *Task {
	for _, t := range d.Tasks {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func TestForwardAfterDependency(t *testing.T) {
	Convey("Given a task that depends on a task defined later", t, func() {
		src := "gantt\ndateFormat YYYY-MM-DD\nTask A : a, after b, 5d\nTask B : b, 2024-01-01, 3d"

		Convey("When parsing", func() {
			d, err := Parse(src)

			Convey("Then the forward dependency resolves and A starts at B's end", func() {
				So(err, ShouldBeNil)
				a := findTask(d, "a")
				b := findTask(d, "b")
				So(a, ShouldNotBeNil)
				So(b, ShouldNotBeNil)
				So(a.Start.IsZero(), ShouldBeFalse)
				So(a.Start.Equal(b.End()), ShouldBeTrue)
			})
		})
	})
}

func TestInvalidStartDateErrors(t *testing.T) {
	Convey("Given a start date that doesn't match dateFormat", t, func() {
		src := "gantt\ndateFormat YYYY\nTask A : a, 2024-01-01, 5d"

		Convey("When parsing", func() {
			_, err := Parse(src)

			Convey("Then it returns an error instead of dropping the task", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

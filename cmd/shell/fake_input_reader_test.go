package shell

import "testing"

func TestFakeInputReader(t *testing.T) {
	f := &FakeInputReader{}

	f.SetReader(nil)

	if !f.CalledSetReader {
		t.Errorf("failed to assert calling method SetReader on FakeInputReader")
	}

	_ = f.GetReader()

	if !f.CalledGetReader {
		t.Errorf("failed to assert calling method GetReader on FakeInputReader")
	}
}

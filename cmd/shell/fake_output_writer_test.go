package shell

import "testing"

func TestFakeOutputWriter(t *testing.T) {
	f := &FakeOutputWriter{}

	f.SetWriter(nil)

	if !f.CalledSetWriter {
		t.Errorf("failed to assert calling method SetWriter on FakeOutputWriter")
	}

	_ = f.GetWriter()

	if !f.CalledGetWriter {
		t.Errorf("failed to assert calling method GetWriter on FakeOutputWriter")
	}

	f.Println()

	if !f.CalledPrintln {
		t.Errorf("failed to assert calling method Println on FakeOutputWriter")
	}

	f.Printf("")

	if !f.CalledPrintf {
		t.Errorf("failed to assert calling method Printf on FakeOutputWriter")
	}

	f.Error(nil)

	if !f.CalledError {
		t.Errorf("failed to assert calling method Error on FakeOutputWriter")
	}

	f.Warning()

	if !f.CalledWarning {
		t.Errorf("failed to assert calling method Warning on FakeOutputWriter")
	}

	f.Success()

	if !f.CalledSuccess {
		t.Errorf("failed to assert calling method Success on FakeOutputWriter")
	}
}

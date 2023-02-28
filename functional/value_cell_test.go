package functional

import "testing"

func TestValueCellWatcher(t *testing.T) {
	t.Parallel()
	vc := CreateValueCell(1)
	var o, n int
	vc.AddWatcher(func(old int, new int) {
		o = old
		n = new
	})
	vc.Set(2)
	if o != 1 || n != 2 {
		t.Fatalf("expected old 1, got %d; expected new 2, got %d", o, n)
	}
	if vc.Get() != 2 {
		t.Fatalf("Expected 2 as result, got %d", vc.Get())
	}
}

func TestValueCellComplex(t *testing.T) {
	t.Parallel()
	vc1 := CreateValueCell(1)
	vc2 := CreateValueCell(2)
	vc3 := CreateValueCell(0)
	w := func(_ int, new int) {
		vc3.Set(vc1.Get() + vc2.Get())
	}
	vc1.AddWatcher(w)
	vc2.AddWatcher(w)
	vc1.Set(2)
	vc2.Set(3)
	if vc3.Get() != 5 {
		t.Fatalf("expected 5 as result, got %d", vc3.Get())
	}
}

func TestValueCellLockCondition(t *testing.T) {
	t.Parallel()
	vc1 := CreateValueCell(1)
	vc1.AddWatcher(func(_ int, _ int) {
		success := vc1.Set(3)
		if success {
			t.Fatalf("expected no success when setting value cell in watcher")
		}
	})
	success := vc1.Set(2)
	if !success {
		t.Fatalf("expected success when setting value cell")
	}
}

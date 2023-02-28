package functional

import "testing"

func TestFormulaCellOneHigher(t *testing.T) {
	t.Parallel()
	vc := CreateValueCell(1)
	fc := CreateFormulaCell(vc, func(old int, new int) int {
		return new + 1
	})
	vc.Set(2)
	if fc.Get() != 3 {
		t.Fatalf("expected 3, received %d", fc.Get())
	}
}

func TestFormulaCellWatcher(t *testing.T) {
	t.Parallel()
	vc := CreateValueCell(1)
	fc := CreateFormulaCell(vc, func(old int, new int) int {
		return new + 1
	})
	var o, n int
	fc.AddWatcher(func(old int, new int) {
		o = old
		n = new
	})
	vc.Set(2)
	if o != 2 || n != 3 {
		t.Fatalf("expected old value 2, got %d; expected new value 3, got %d", o, n)
	}
}

package booter

import "testing"

func TestBooter_RegisterGetWithDep(t *testing.T) {
	booter := NewBooter()

	booter.Register("test", func() (interface{}, error) {
		dep := booter.Get("test.dep")
		return "test-service-instance with " + dep.(string), nil
	})
	booter.Register("test.dep", func() (interface{}, error) {
		return "dep", nil
	})

	svcInstance := booter.Get("test")
	if svcInstance != "test-service-instance with dep" {
		t.Error("svcInstance != `test-service-instance`")
	}
}

func TestBooter_GetOnce(t *testing.T) {

	svcRegistraCalled := 0

	booter := NewBooter()

	booter.Register("test", func() (interface{}, error) {
		svcRegistraCalled++
		return "test-service-instance", nil
	})

	_ = booter.Get("test")
	_ = booter.Get("test")
	_ = booter.Get("test")

	if svcRegistraCalled != 1 {
		t.Error("svcRegistraCalled != 1")
	}
}

func TestBooter_CircularDependency(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Get did not panic")
		}
	}()

	booter := NewBooter()
	booter.Register("svc1", func() (interface{}, error) {
		_ = booter.Get("svc2")
		return "svc1", nil
	})
	booter.Register("svc2", func() (interface{}, error) {
		_ = booter.Get("svc3")
		return "svc2", nil
	})

	booter.Register("svc3", func() (interface{}, error) {
		_ = booter.Get("svc1")
		return "svc3", nil
	})
	_ = booter.Get("svc1")
}

func TestBooter_MustGet(t *testing.T) {
	booter := NewBooter()

	booter.Register("test", func() (interface{}, error) {
		return "test-service-instance", nil
	})

	svcInstance := booter.Get("test")

	if svcInstance != "test-service-instance" {
		t.Error("svcInstance != `test-service-instance`")
	}
}

func TestBooter_MustGetPanic(t *testing.T) {
	booter := NewBooter()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet did not panic")
		}
	}()

	booter.Get("not-registered")
}

func TestBooter_Cache(t *testing.T) {
	booter := NewBooterWithCached(map[string]interface{}{
		"foo": "bar",
	})

	booter.Cache("test", "test-service-instance")

	if booter.Get("test").(string) != "test-service-instance" {
		t.Error("test service is not `test-service-instance`")
	}

	if booter.Get("foo").(string) != "bar" {
		t.Error("foo service is not `bar`")
	}
}

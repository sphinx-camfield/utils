package booter

import "testing"

func TestBooter_RegisterGetWithDep(t *testing.T) {
	booter := NewBooter()

	booter.Register("test", func() (interface{}, error) {
		dep, err := booter.Get("test.dep")

		if err != nil {
			return nil, err
		}

		return "test-service-instance with " + dep.(string), nil
	})

	booter.Register("test.dep", func() (interface{}, error) {
		return "dep", nil
	})

	svcInstance, err := booter.Get("test")

	if err != nil {
		t.Error(err)
	}

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

	_, _ = booter.Get("test")
	_, _ = booter.Get("test")
	_, _ = booter.Get("test")

	if svcRegistraCalled != 1 {
		t.Error("svcRegistraCalled != 1")
	}
}

func TestBooter_CircularDependency(t *testing.T) {

	booter := NewBooter()

	booter.Register("svc1", func() (interface{}, error) {
		_, err := booter.Get("svc2")
		if err != nil {
			return nil, err
		}
		return "svc1", nil
	})

	booter.Register("svc2", func() (interface{}, error) {
		_, err := booter.Get("svc3")
		if err != nil {
			return nil, err
		}
		return "svc2", nil
	})

	booter.Register("svc3", func() (interface{}, error) {
		_, err := booter.Get("svc1")
		if err != nil {
			return nil, err
		}
		return "svc3", nil
	})

	_, err := booter.Get("svc1")

	if err == nil {
		t.Error("err == nil")
	}

	if err.Error() != "circular dependency detected: [svc1 svc2 svc3 svc1]" {

		t.Error("err.Error() != `circular dependency detected: [svc1 svc2 svc3 svc1]`")
	}
}

func TestBooter_MustGet(t *testing.T) {
	booter := NewBooter()

	booter.Register("test", func() (interface{}, error) {
		return "test-service-instance", nil
	})

	svcInstance := booter.MustGet("test")

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

	booter.MustGet("not-registered")
}

func TestBooter_Cache(t *testing.T) {
	booter := NewBooterWithCached(map[string]interface{}{
		"foo": "bar",
	})

	booter.Cache("test", "test-service-instance")

	if booter.MustGet("test").(string) != "test-service-instance" {
		t.Error("test service is not `test-service-instance`")
	}

	if booter.MustGet("foo").(string) != "bar" {
		t.Error("foo service is not `bar`")
	}
}

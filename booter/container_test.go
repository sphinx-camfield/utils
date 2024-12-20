package booter

import "testing"

func TestBooter_RegisterGetWithDep(t *testing.T) {
	booter := NewContainer()

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

	booter := NewContainer()

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

	booter := NewContainer()
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
	booter := NewContainer()

	booter.Register("test", func() (interface{}, error) {
		return "test-service-instance", nil
	})

	svcInstance := booter.Get("test")

	if svcInstance != "test-service-instance" {
		t.Error("svcInstance != `test-service-instance`")
	}
}

func TestBooter_MustGetPanic(t *testing.T) {
	booter := NewContainer()

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

func TestBooter_Alias(t *testing.T) {
	booter := NewContainer()

	booter.Register("test", func() (interface{}, error) {
		return "test-service-instance", nil
	})

	booter.Alias("test", "test-alias")

	if booter.Get("test-alias").(string) != "test-service-instance" {
		t.Error("test-alias service is not `test-service-instance`")
	}
}

func TestBooter_AliasSameAsSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Alias did not panic")
		}
	}()

	booter := NewContainer()
	booter.Alias("test", "test")
}

func TestBooter_CyclicAliasToAlias(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Alias did not panic")
		}
	}()
	booter := NewContainer()
	booter.Alias("test", "alias")
	booter.Alias("alias", "test")
}

func TestBooter_AliasToAlias(t *testing.T) {
	booter := NewContainer()
	booter.Register("test", func() (interface{}, error) {
		return "test-service-instance", nil
	})
	booter.Alias("test", "alias")
	booter.Alias("alias", "alias2")
	booter.Alias("test", "alias3")

	if booter.Get("alias2").(string) != "test-service-instance" {
		t.Error("alias2 service is not `test-service-instance`")
	}

	if booter.Get("alias3").(string) != "test-service-instance" {
		t.Error("alias3 service is not `test-service-instance`")
	}
}

func TestBooter_AliasToNonExistence(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Alias did not panic")
		}
	}()
	booter := NewContainer()
	booter.Alias("test", "alias")
	booter.Get("alias")
}

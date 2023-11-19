# Floaty
☁️ Prevent loopbacks... Without traceability.

This is a Caddy plugin that adds customizable rolling instance IDs, permitting loopback prevention while maximizing difficulty in tracing.

Schema:

```sh
http://example.com:8080 {
	floaty [serviceName] {
		[fieldName [idLength [rollDuration]]]
	}
}
```

Documentation available at [kb.ltgc.cc](https://kb.ltgc.cc/floaty/).
package nvppayload

type (
	// Serializer interface for payloads that can be serialized
	Serializer interface {
		Serialize() (string, error)
	}
)

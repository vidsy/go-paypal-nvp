package payload

type (
	// Serializer interface for payloads that can be serialized.
	Serializer interface {
		Serialize() (string, error)
		SetCredentials(string, string, string, string)
	}
)

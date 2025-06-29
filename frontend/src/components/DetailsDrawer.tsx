interface DetailsDrawerProps<TData> {
  data: TData;
  fields: { key: keyof TData; label: string }[];
}

export function DetailsDrawer<TData>({ data, fields }: DetailsDrawerProps<TData>) {
  return (
    <div className="p-4">
      <h3 className="text-lg font-semibold mb-2">Details</h3>
      {
        fields.map((field) => (
          <p key={String(field.key)}>
            <strong>{field.label}:</strong> {String(data[field.key])}
          </p>
        ))
      }
    </div>
  );
}
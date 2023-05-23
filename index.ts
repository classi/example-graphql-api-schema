import { basename, join } from "node:path";
import schemaFiles from "./schema-files.json";

export interface SchemaFile {
  readonly name: string;
  readonly path: string;
}

export interface SchemaFilesAggregate {
  [name: string]: SchemaFile;
}

const files = schemaFiles.map(
  (fp): SchemaFile => ({
    name: basename(fp),
    path: join(__dirname, fp),
  })
);

const filesByName = files.reduce<SchemaFilesAggregate>(
  (a, f) => ({
    ...a,
    [f.name]: f,
  }),
  {}
);

export { files, filesByName };

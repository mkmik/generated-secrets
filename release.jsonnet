(import 'controller.jsonnet') {
  fields+: {
    appImage: importstr 'image.txt',
  },
}

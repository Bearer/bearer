import attach from './attachConfig'

test('attaching the config', () => {
  const archive = {
    append: jest.fn()
  }
  const content = '{"integration_uuid":"uuid","functions":["get","put"]}'
  const fileName = { name: 'bearer.config.json' }

  attach(archive, content, fileName)
  expect(archive.append).toHaveBeenCalledWith(content, fileName)
})

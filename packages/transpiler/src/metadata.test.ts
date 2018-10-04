import Metadata from './metadata'

describe('Metadata', () => {
  describe('constructor', () => {
    it(' default to undefined prefix and suffix', () => {
      const meta = new Metadata()
      expect(meta.prefix).toBeFalsy()
      expect(meta.suffix).toBeFalsy()
    })

    it('set suffix and prefix', () => {
      const meta = new Metadata('prefix', 'suffix')
      expect(meta.prefix).toEqual('prefix')
      expect(meta.suffix).toEqual('suffix')
    })
  })

  describe('registerComponent', () => {
    // it does not play well with updated component

    it('add component if does not exits and sort components', () => {
      const meta = new Metadata()
      meta.registerComponent({
        classname: 'BClass',
        isRoot: true,
        initialTagName: 'b-class',
        finalTagName: 'bearer-b-class'
      })

      meta.registerComponent({
        classname: 'AClass',
        isRoot: true,
        initialTagName: 'a-class',
        finalTagName: 'bearer-a-class'
      })

      expect(meta.components).toEqual([
        {
          classname: 'AClass',
          isRoot: true,
          initialTagName: 'a-class',
          finalTagName: 'bearer-a-class'
        },
        {
          classname: 'BClass',
          isRoot: true,
          initialTagName: 'b-class',
          finalTagName: 'bearer-b-class'
        }
      ])
    })

    it('replaces existing component', () => {
      const meta = new Metadata()
      meta.registerComponent({
        classname: 'BClass',
        isRoot: true,
        initialTagName: 'b-class',
        finalTagName: 'bearer-b-class'
      })

      meta.registerComponent({
        classname: 'UpdatedClass',
        isRoot: true,
        initialTagName: 'b-class',
        finalTagName: 'bearer-b-class'
      })

      expect(meta.components).toEqual([
        {
          classname: 'UpdatedClass',
          isRoot: true,
          initialTagName: 'b-class',
          finalTagName: 'bearer-b-class'
        }
      ])
    })
  })
})

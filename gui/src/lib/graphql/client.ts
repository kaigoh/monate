import { HoudiniClient } from '$houdini';

export default new HoudiniClient({
  url: '/query',
  throwOnError: {
    operations: ['mutation']
  }
});

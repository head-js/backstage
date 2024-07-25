<template>
  <ElAside width="180px">
    <ElMenu :router="true">
      <ElSubmenu index="1">
        <template slot="title">
          <i class="el-icon-location" />
          <span>导航一</span>
        </template>
        <ElMenuItemGroup>
          <template slot="title">分组一</template>
          <ElMenuItem index="1-1">选项一</ElMenuItem>
          <ElMenuItem index="1-2">选项二</ElMenuItem>
        </ElMenuItemGroup>
      </ElSubmenu>
      <ElMenuItem index="/account">
        <i class="el-icon-setting" />
        <span>Account</span>
      </ElMenuItem>
      <ElMenuItem index="/options">
        <i class="el-icon-setting" />
        <span>Options</span>
      </ElMenuItem>
    </ElMenu>
  </ElAside>
</template>

<script>
import { Aside as ElAside, Menu as ElMenu, Submenu as ElSubmenu, MenuItemGroup as ElMenuItemGroup, MenuItem as ElMenuItem } from 'element-ui';

const namespace = 'aside-menu';

export default {
  name: 'AsideMenu',

  components: { ElAside, ElMenu, ElSubmenu, ElMenuItemGroup, ElMenuItem },

  computed: {},
};
</script>

<style lang="css" scoped>
.el-aside {
  height: 90vh;
}

.el-aside .el-submenu .el-menu-item {
  min-width: 180px;
}
</style>

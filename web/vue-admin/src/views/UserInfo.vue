<template>
  <section>
    <!--工具条-->
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.name" placeholder="输入用户ID或姓名查询"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="getUsers">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </el-col>

    <!--列表-->
    <el-table :data="datas" highlight-current-row v-loading="listLoading" @selection-change="selsChange"
              style="width: 100%;">
      <el-table-column type="selection" width="55">
      </el-table-column>
      <el-table-column type="index" width="60">
      </el-table-column>
      <el-table-column prop="uid" label="用户" width="120" sortable>
      </el-table-column>
      <el-table-column prop="uname" label="姓名" width="100" sortable>
      </el-table-column>
      <el-table-column prop="department" label="业务部" width="100" sortable>
      </el-table-column>
      <el-table-column prop="group" label="组" width="120" sortable>
      </el-table-column>
      <el-table-column prop="comment" label="备注" min-width="180" sortable>
      </el-table-column>
      <el-table-column prop="enabled" label="状态" min-width="180" :formatter="formatEnabled" sortable>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template scope="scope">
          <el-button size="small" @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDel(scope.$index, scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!--工具条-->
    <el-col :span="24" class="toolbar">
      <el-button type="danger" @click="batchRemove" :disabled="this.sels.length===0">批量删除</el-button>
      <el-pagination layout="prev, pager, next" @current-change="handleCurrentChange" :page-size="20" :total="total"
                     style="float:right;">
      </el-pagination>
    </el-col>

    <!--编辑界面-->
    <el-dialog title="编辑" :visible.sync="editUserFormVisible" :close-on-click-modal="false">
      <el-form :model="editUserForm" label-width="80px" :rules="UserFormRules" ref="editUserForm">
        <el-form-item label="用户ID" prop="uid">
          <el-input v-model="editUserForm.uid" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="uname">
          <el-input v-model="editUserForm.uname"></el-input>
        </el-form-item>
        <el-form-item label="业务部" prop="department">
          <el-input v-model="editUserForm.department"></el-input>
        </el-form-item>
        <el-form-item label="组" prop="group">
          <el-input v-model="editUserForm.group"></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="editUserForm.comment"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="editUserForm.enabled">
            <el-radio class="radio" :label=true>启用</el-radio>
            <el-radio class="radio" :label=false>未启用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editUserFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
      </div>
    </el-dialog>

    <!--新增界面-->
    <el-dialog title="新增" :visible.sync="addUserFormVisible" :close-on-click-modal="false">
      <el-form :model="addUserForm" label-width="80px" :rules="UserFormRules" ref="addUserForm">
        <el-form-item label="用户ID" prop="uid">
          <el-input v-model="addUserForm.uid" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="uname">
          <el-input v-model="addUserForm.uname" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="业务部" prop="department">
          <el-input v-model="addUserForm.department" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="组" prop="group">
          <el-input v-model="addUserForm.group" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="addUserForm.comment"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="addUserForm.enabled">
            <el-radio class="radio" :label=true>启用</el-radio>
            <el-radio class="radio" :label=false>未启用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addUserFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
import util from '../common/js/util'
import { getUserListPage, removeUser, batchRemoveUser, editUser, addUser } from '../api/api'

export default {
  data () {
    return {
      filters: {
        name: ''
      },
      datas: [],
      total: 0,
      page: 1,
      listLoading: false,
      sels: [], // 列表选中列

      editUserFormVisible: false, // 编辑界面是否显示
      editLoading: false,
      UserFormRules: {
        uid: [
          { required: true, message: '请输入用户ID', trigger: 'blur' }
        ],
        uname: [
          { required: true, message: '请输入用户姓名', trigger: 'blur' }
        ],
        department: [
          { required: true, message: '请输入用户所属部门', trigger: 'blur' }
        ],
        group: [
          { required: true, message: '请输入用户所属组', trigger: 'blur' }
        ]
      },
      // 编辑界面数据
      editUserForm: {
        id: 0,
        uid: '',
        uname: '',
        department: '',
        group: '',
        comment: '',
        enabled: true

      },

      addUserFormVisible: false, // 新增界面是否显示
      addLoading: false,
      addUserFormRules: {
        name: [
          { required: true, message: '请输入姓名', trigger: 'blur' }
        ]
      },
      // 新增界面数据
      addUserForm: {
        uid: '',
        uname: '',
        department: '',
        group: '',
        comment: '',
        enabled: true
      }

    }
  },
  methods: {
    // 状态显示转换
    formatEnabled: function (row, column) {
      return row.enabled ? '启用' : '未启用'
    },
    handleCurrentChange (val) {
      this.page = val
      this.getUsers()
    },
    // 获取用户列表
    getUsers () {
      let para = {
        page: this.page,
        name: this.filters.name
      }
      this.listLoading = true
      getUserListPage(para).then((res) => {
        this.total = res.data.total
        this.datas = JSON.parse(res.data.data)
        this.listLoading = false
      })
    },
    // 删除
    handleDel: function (index, row) {
      this.$confirm('确认删除该记录吗?', '提示', {
        type: 'warning'
      }).then(() => {
        this.listLoading = true
        let para = { id: row.id }
        removeUser(para).then((res) => {
          this.listLoading = false
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.getUsers()
        })
      }).catch(() => {

      })
    },
    // 显示编辑界面
    handleEdit: function (index, row) {
      this.editUserFormVisible = true
      this.editUserForm = Object.assign({}, row)
    },
    // 显示新增界面
    handleAdd: function () {
      this.addUserFormVisible = true
      this.addUserForm = {
        uid: '',
        uname: '',
        department: '',
        group: '',
        comment: '',
        enabled: true
      }
    },
    // 编辑
    editSubmit: function () {
      this.$refs.editUserForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认提交吗？', '提示', {}).then(() => {
            this.editLoading = true
            let para = Object.assign({}, this.editUserForm)
            editUser(para).then((res) => {
              this.editLoading = false
              let { code, msg, detail } = res
              if (code !== 200) {
                this.$message({
                  message: '提交失败',
                  type: 'error'
                })
              } else {
                this.$message({
                  message: '提交成功',
                  type: 'success'
                })
                this.$refs['editUserForm'].resetFields()
                this.editUserFormVisible = false
                this.getUsers()
              }
            })
          })
        }
      })
    },
    // 新增
    addSubmit: function () {
      this.$refs.addUserForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认提交吗？', '提示', {}).then(() => {
            this.addLoading = true
            // NProgress.start();
            let para = Object.assign({}, this.addUserForm)
            addUser(para).then((res) => {
              this.addLoading = false
              // NProgress.done();
              let { code, msg, detail } = res
              if (code !== 200) {
                this.$message({
                  message: '提交失败',
                  type: 'error'
                })
              } else {
                this.$message({
                  message: '提交成功',
                  type: 'success'
                })
                this.$refs['addUserForm'].resetFields()
                this.addUserFormVisible = false
                this.getUsers()
              }
            })
          })
        }
      })
    },
    selsChange: function (sels) {
      this.sels = sels
    },
    // 批量删除
    batchRemove: function () {
      var ids = this.sels.map(item => item.id)
      this.$confirm('确认删除选中记录吗？', '提示', {
        type: 'warning'
      }).then(() => {
        this.listLoading = true
        let para = { ids: ids }
        batchRemoveUser(para).then((res) => {
          this.listLoading = false
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.getUsers()
        })
      }).catch(() => {

      })
    }
  },
  mounted () {
    this.getUsers()
  }
}

</script>

<style scoped>

</style>

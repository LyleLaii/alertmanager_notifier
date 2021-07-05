<template>
  <section>
    <!--工具条-->
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.name" placeholder="输入用户ID或姓名查询"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="getRotas">查询</el-button>
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
      <el-table-column prop="uname" label="姓名" width="100"  sortable>
      </el-table-column>
      <el-table-column prop="department" label="业务部" width="100" sortable>
      </el-table-column>
      <el-table-column prop="group" label="组" min-width="180" sortable>
      </el-table-column>
      <el-table-column prop="date" label="值班日期" width="120" sortable>
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
    <el-dialog title="编辑" :visible.sync="editRotaFormVisible" :close-on-click-modal="false">
      <el-form :model="editRotaForm" label-width="80px" :rules="editRotaFormRules" ref="editRotaForm">
        <el-form-item label="用户" prop="uid">
          <el-input v-model="editRotaForm.uid" auto-complete="off" disabled></el-input>
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="editRotaForm.uname" disabled></el-input>
        </el-form-item>
        <el-form-item label="业务部">
          <el-input v-model="editRotaForm.department" disabled></el-input>
        </el-form-item>
        <el-form-item label="组">
          <el-input v-model="editRotaForm.group" disabled></el-input>
        </el-form-item>
        <el-form-item label="值班日期">
          <el-date-picker type="date" placeholder="选择日期" v-model="editRotaForm.date"></el-date-picker>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editRotaFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
      </div>
    </el-dialog>

    <!--新增界面-->
    <el-dialog title="新增" :visible.sync="addRotaFormVisible" :close-on-click-modal="false">
      <el-form :model="addRotaForm" label-width="80px" :rules="addRotaFormRules" ref="addRotaForm">
        <el-form-item label="用户" prop="uid">
          <el-select v-model="addRotaForm.uid" placeholder="输入用户ID" @change="getSelected" filterable clearable>
          <el-option v-for="item in allusers" :label="item" :value="item" v-bind:key="item"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="addRotaForm.uname" disabled></el-input>
        </el-form-item>
        <el-form-item label="业务部">
          <el-input v-model="addRotaForm.department" disabled></el-input>
        </el-form-item>
        <el-form-item label="组">
          <el-input v-model="addRotaForm.group" disabled></el-input>
        </el-form-item>
        <el-form-item label="值班日期">
          <el-date-picker type="date" placeholder="选择日期" v-model="addRotaForm.date"></el-date-picker>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addRotaFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
import util from '../common/js/util'
import { getRotaListPage, removeRota, batchRemoveRota, editRota, addRota, getAllUsers, getUserInfo } from '../api/api'

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
      allusers: [],

      editRotaFormVisible: false, // 编辑界面是否显示
      editLoading: false,
      editRotaFormRules: {
        uid: [
          { required: true, message: '请选择用户', trigger: 'blur' }
        ],
        date: [
          { required: true, message: '请设定值班时间', trigger: 'blur' }
        ]
      },
      // 编辑界面数据
      editRotaForm: {
        id: 0,
        uid: '',
        uname: '',
        department: '',
        group: '',
        date: ''
      },

      addRotaFormVisible: false, // 新增界面是否显示
      addLoading: false,
      addRotaFormRules: {
        uid: [
          { required: true, message: '请选择用户', trigger: 'blur' }
        ],
        date: [
          { required: true, message: '请设定值班时间', trigger: 'blur' }
        ]
      },
      // 新增界面数据
      addRotaForm: {
        uid: '',
        uname: '',
        department: '',
        group: '',
        date: ''
      }

    }
  },
  methods: {
    handleCurrentChange (val) {
      this.page = val
      this.getRotas()
    },
    getUser () {
      getAllUsers().then(res => {
        this.allusers = res.data.data
      }).catch(err => {
        console.log(err)
      })
    },

    getSelected (val) {
      let para = {
        uid: val
      }
      getUserInfo(para).then(res => {
        this.addRotaForm = Object.assign({}, res.data.data)
        // this.addRotaForm.uname = data.name
        // this.addRotaForm.department = data.department
        // this.addRotaForm.group = data.group
      })
    },
    // 获取用户列表
    getRotas () {
      let para = {
        page: this.page,
        name: this.filters.name
      }
      this.listLoading = true
      getRotaListPage(para).then((res) => {
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
        removeRota(para).then((res) => {
          this.listLoading = false
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.getRotas()
        })
      }).catch(() => {

      })
    },
    // 显示编辑界面
    handleEdit: function (index, row) {
      this.editRotaFormVisible = true
      this.editRotaForm = Object.assign({}, row)
    },
    // 显示新增界面
    handleAdd: function () {
      this.addRotaFormVisible = true
      this.addRotaForm = {
        uid: '',
        uname: '',
        department: '',
        group: '',
        date: ''
      }
    },
    // 编辑
    editSubmit: function () {
      this.$refs.editRotaForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认提交吗？', '提示', {}).then(() => {
            this.editLoading = true
            // let para = Object.assign({}, this.editRotaForm)
            let para = { id: this.editRotaForm.id,
              uid: this.editRotaForm.uid,
              date: this.editRotaForm.date }
            para.date = (!para.date || para.date === '') ? '' : util.formatDate.format(new Date(para.date), 'yyyy-MM-dd')
            editRota(para).then((res) => {
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
                this.$refs['editRotaForm'].resetFields()
                this.editRotaFormVisible = false
                this.getRotas()
              }
            })
          })
        }
      })
    },
    // 新增
    addSubmit: function () {
      this.$refs.addRotaForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认提交吗？', '提示', {}).then(() => {
            this.addLoading = true
            // NProgress.start();
            // let para = Object.assign({}, this.addRotaForm)
            let para = { uid: this.addRotaForm.uid,
              date: this.addRotaForm.date }
            para.date = (!para.date || para.date === '') ? '' : util.formatDate.format(new Date(para.date), 'yyyy-MM-dd')
            addRota(para).then((res) => {
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
                this.$refs['addRotaForm'].resetFields()
                this.addRotaFormVisible = false
                this.getRotas()
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
        batchRemoveRota(para).then((res) => {
          this.listLoading = false
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.getRotas()
        })
      }).catch(() => {

      })
    }
  },
  mounted () {
    this.getRotas()
    this.getUser()
  }
}

</script>

<style scoped>

</style>

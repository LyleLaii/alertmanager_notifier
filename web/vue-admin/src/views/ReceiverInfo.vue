<template>
  <section>
    <!--工具条-->
    <el-col :span="24" class="toolbar" style="padding-bottom: 0px;">
      <el-form :inline="true" :model="filters">
        <el-form-item>
          <el-input v-model="filters.name" placeholder="输入用户ID或姓名查询"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" v-on:click="getReceiverInfos">查询</el-button>
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
      <el-table-column prop="uname" label="姓名" width="120" sortable>
      </el-table-column>
      <el-table-column prop="receiverType" label="接收方式" width="100"  sortable>
      </el-table-column>
      <el-table-column prop="receiverName" label="接收名称" width="100" sortable>
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
    <el-dialog title="编辑" :visible.sync="editReceiverInfoFormVisible" :close-on-click-modal="false">
      <el-form :model="editReceiverInfoForm" label-width="80px" :rules="ReceiverInfoFormRules" ref="editReceiverInfoForm">
        <el-form-item label="用户" prop="uid">
          <el-input v-model="editReceiverInfoForm.uid" auto-complete="off" disabled></el-input>
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="editReceiverInfoForm.uname" disabled></el-input>
        </el-form-item>
        <el-form-item label="接收方式" prop="receiverType">
          <el-input v-model="editReceiverInfoForm.receiverType" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="接收名称" prop="receiverName">
          <el-input v-model="editReceiverInfoForm.receiverName" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="editReceiverInfoFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="editSubmit" :loading="editLoading">提交</el-button>
      </div>
    </el-dialog>

    <!--新增界面-->
    <el-dialog title="新增" :visible.sync="addReceiverInfoFormVisible" :close-on-click-modal="false">
      <el-form :model="addReceiverInfoForm" label-width="80px" :rules="ReceiverInfoFormRules" ref="addReceiverInfoForm">
        <el-form-item label="用户" prop="uid">
          <el-select v-model="addReceiverInfoForm.uid" placeholder="输入用户ID" @change="getSelected" filterable clearable>
            <el-option v-for="item in allusers" :label="item" :value="item" v-bind:key="item"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="姓名" prop="uname">
          <el-input v-model="addReceiverInfoForm.uname" disabled></el-input>
        </el-form-item>
        <el-form-item label="接收方式" prop="receiverType">
          <el-input v-model="addReceiverInfoForm.receiverType" auto-complete="off"></el-input>
        </el-form-item>
        <el-form-item label="接收名称" prop="receiverName">
          <el-input v-model="addReceiverInfoForm.receiverName" auto-complete="off"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click.native="addReceiverInfoFormVisible = false">取消</el-button>
        <el-button type="primary" @click.native="addSubmit" :loading="addLoading">提交</el-button>
      </div>
    </el-dialog>
  </section>
</template>

<script>
import util from '../common/js/util'
import { getReceiverInfoListPage, removeReceiverInfo, batchRemoveReceiverInfo, editReceiverInfo, addReceiverInfo, getAllUsers, getUserInfo } from '../api/api'

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

      editReceiverInfoFormVisible: false, // 编辑界面是否显示
      editLoading: false,
      ReceiverInfoFormRules: {
        uid: [
          { required: true, message: '请选择用户', trigger: 'blur' }
        ],
        receiverType: [
          { required: true, message: '请填写接收类型', trigger: 'blur' }
        ],
        receiverName: [
          { required: true, message: '请填写接收名称', trigger: 'blur' }
        ]
      },
      // 编辑界面数据
      editReceiverInfoForm: {
        id: 0,
        uid: '',
        uname: '',
        receiverType: '',
        receiverName: ''
      },

      addReceiverInfoFormVisible: false, // 新增界面是否显示
      addLoading: false,
      addReceiverInfoFormRules: {
        uid: [
          { required: true, message: '请选择用户', trigger: 'blur' }
        ],
        receiverType: [
          { required: true, message: '请填写接收类型', trigger: 'blur' }
        ],
        receiverName: [
          { required: true, message: '请填写接收名称', trigger: 'blur' }
        ]
      },
      // 新增界面数据
      addReceiverInfoForm: {
        uid: '',
        uname: '',
        receiverType: '',
        receiverName: ''
      }

    }
  },
  methods: {
    handleCurrentChange (val) {
      this.page = val
      this.getReceiverInfos()
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
        this.addReceiverInfoForm = Object.assign({}, res.data.data)
      })
    },
    // 获取用户列表
    getReceiverInfos () {
      let para = {
        page: this.page,
        name: this.filters.name
      }
      this.listLoading = true
      getReceiverInfoListPage(para).then((res) => {
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
        removeReceiverInfo(para).then((res) => {
          this.listLoading = false
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.getReceiverInfos()
        })
      }).catch(() => {

      })
    },
    // 显示编辑界面
    handleEdit: function (index, row) {
      this.editReceiverInfoFormVisible = true
      this.editReceiverInfoForm = Object.assign({}, row)
    },
    // 显示新增界面
    handleAdd: function () {
      this.addReceiverInfoFormVisible = true
      this.addReceiverInfoForm = {
        uid: '',
        uname: '',
        receiverType: '',
        receiverName: ''
      }
    },
    // 编辑
    editSubmit: function () {
      this.$refs.editReceiverInfoForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认提交吗？', '提示', {}).then(() => {
            this.editLoading = true
            // let para = Object.assign({}, this.editReceiverInfoForm)
            let para = { id: this.editReceiverInfoForm.id,
              uid: this.editReceiverInfoForm.uid,
              receiverType: this.editReceiverInfoForm.receiverType,
              receiverName: this.editReceiverInfoForm.receiverName }
            editReceiverInfo(para).then((res) => {
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
                this.$refs['editReceiverInfoForm'].resetFields()
                this.editReceiverInfoFormVisible = false
                this.getReceiverInfos()
              }
            })
          })
        }
      })
    },
    // 新增
    addSubmit: function () {
      this.$refs.addReceiverInfoForm.validate((valid) => {
        if (valid) {
          this.$confirm('确认提交吗？', '提示', {}).then(() => {
            this.addLoading = true
            // NProgress.start();
            // let para = Object.assign({}, this.addReceiverInfoForm)
            let para = { uid: this.addReceiverInfoForm.uid,
              receiverType: this.addReceiverInfoForm.receiverType,
              receiverName: this.addReceiverInfoForm.receiverName }
            addReceiverInfo(para).then((res) => {
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
                this.$refs['addReceiverInfoForm'].resetFields()
                this.addReceiverInfoFormVisible = false
                this.getReceiverInfos()
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
        batchRemoveReceiverInfo(para).then((res) => {
          this.listLoading = false
          this.$message({
            message: '删除成功',
            type: 'success'
          })
          this.getReceiverInfos()
        })
      }).catch(() => {

      })
    }
  },
  mounted () {
    this.getReceiverInfos()
    this.getUser()
  }
}

</script>

<style scoped>

</style>

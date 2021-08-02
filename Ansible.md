### register、debug变量
```
---
- hosts: all
  gather_facts: no
  tasks:
    - name: register vars
      shell: hostname
      register: info
    - name: display vars
      debug: msg="{{info.stdout}}"
```
**register**表示第一个shell执行完后，将得到的结果存到info里<br>
**info**是一个key value字典<br>
**debug**表示输出info.stdout的具体内容<br>


### became功能
您可以设置控制become进入play或task级别的指令。您可以通过设置连接变量来覆盖这些变量，而连接变量通常在一台主机之间会有所不同。这些变量和指令是独立的。比如，可以单独设置become_user，而不设置become。


|become | 是否进行提权操作。如果需要，设置为yes|    是|
| :----  | ----  | ---  |
|become_user | root 设置为具有所需特权的用户-您想要成为的用户，而不是您登录时使用的用户 |否|
|become_method | sudo  权限工具，如sudo，su，pfexec，doas，pbrun，dzdo，ksu，runas，machinectl | 否
|become_flags  | play或task级别上，允许为任务或角色使用特定的标志。一种常见的用法是，当shell设置为no login时，将用户更改为nobody。此指令是在Ansible 2.2中添加|

### cp模块
backup参数 : 当远程主机的目标路径中已经存在同名文件，并且与ansible主机中的文件内容不同时，是否对远程主机的文件进行备份，可选值有yes和no，当设置为yes时，会先备份远程主机中的文件，然后再将ansible主机中的文件拷贝到远程主机。
```
 copy: src={{ item }} dest=/etc backup=yes
  with_items:
    - resolv.conf
    - sysctl.conf
    - profile
```

### vars变量 ：使用自定义主机的内置参数
```
vim /etc/ansible/hosts
[webservers]
192.168.1.[31:32]
[webservers:vars]
ansible_ssh_user='root'
ansible_ssh_pass='redhat'  #这个Redhat是自己定义的参数
ansible_ssh_port='22'
```
### 子组分类变量：children
```
vim /etc/ansible/hosts 
[nginx]
192.168.1.31
[apache]
192.168.1.32
[webservers:children]
apache
nginx
[webservers:vars]
ansible_ssh_user='root'
ansible_ssh_pass='redhat'
ansible_ssh_port='22'
```
### {{}}
```
- name: 02 {{ server }}:解压软件包
  unarchive: src=/etc/ansible/roles/deploy/package/{{ env }}/{{ server }}.tar.gz dest={{ PACK_DIR }} owner=holo_usr group=holo_usr
```
{{}}表引用变量  server外部传参数，ansible -e server

### lineinfile模块：针对文件的操作的模块
`lineinfile: path=/etc/cron.allow line="holo_usr"  ` <br>
表示在=/etc/cron.allow中查找holo_user这个字符串，存在则不执行任何操作，不存在则在文本末尾插入

### with_items
介绍：with_items可以用于迭代一个列表或字典，通过{{ item }}获取每次迭代的值。
```
- name: 01 删除用户
  user: name={{ item.name }} group={{ item.group }} state=absent remove=yes
  with_items:
    - { name: 'service' , group: 'servicegroup' }
    - { name: 'opsadmin' , group: 'admingroup' }
    - { name: 'gluster' , group: 'gluster' }
```

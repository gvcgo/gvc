if exists('g:vscode')
	set nobackup
else
	set mouse=a
endif

let mapleader="\<space>"

function! Cond(Cond, ...)
  let opts = get(a:000, 0, {})
  return a:Cond ? opts : extend(opts, { 'on': [], 'for': [] })
endfunction
call plug#begin('~/.config/nvim/plugged')
Plug 'easymotion/vim-easymotion', Cond(!exists('g:vscode'))
Plug 'asvetliakov/vim-easymotion', Cond(exists('g:vscode'), { 'as': 'vsc-easymotion' })
" Plug 'fatih/vim-go'
call plug#end()

noremap H ^
noremap L $

" easymotion相关配置
let g:EasyMotion_smartcase = 0
" easymotion前缀 leader leader
map <Leader> <Plug>(easymotion-prefix)
" 其他键位绑定
map <Leader>l <Plug>(easymotion-lineforward)
map <Leader>j <Plug>(easymotion-j)
map <Leader>k <Plug>(easymotion-k)
map <Leader>h <Plug>(easymotion-linebackward)

if exists('g:vscode')
	" 使用vscode的undo替换nvim的undo
	nnoremap u <Cmd>call VSCodeNotify('undo')<CR>
	" 使用vscode的调试
	nnoremap ge <Cmd>call VSCodeNotify('workbench.action.debug.start', {'when': "!inDebugMode"})<CR>
	nnoremap gr <Cmd>call VSCodeNotify('workbench.action.debug.restart', {'when': "inDebugMode"})<CR>
	nnoremap gs <Cmd>call VSCodeNotify('workbench.action.debug.stop', {'when': "inDebugMode"})<CR>
	nnoremap ga <Cmd>call VSCodeNotify('workbench.debug.action.focusVariablesView', {'when': "inDebugMode"})<CR>
	nnoremap go <Cmd>call VSCodeNotify('workbench.action.debug.stepOut', {'when': "debugState == 'running'"})<CR>
	nnoremap gv <Cmd>call VSCodeNotify('workbench.action.debug.stepOver', {'when': "debugState == 'running'"})<CR>
	nnoremap gn <Cmd>call VSCodeNotify('workbench.action.debug.continue', {'when': "debugState == 'running'"})<CR>
	nnoremap gi <Cmd>call VSCodeNotify('workbench.action.debug.stepInto', {'when': "debugState == 'running'"})<CR>
	nnoremap gp <Cmd>call VSCodeNotify('workbench.action.debug.pause', {'when': "debugState == 'running'"})<CR>
	nnoremap gk <Cmd>call VSCodeNotify('editor.debug.action.toggleBreakpoint')<CR>
	nnoremap gh <Cmd>call VSCodeNotify('editor.debug.action.showDebugHover', {'when': "editorTextFocus && inDebugMode"})<CR>
	" 跳转到下次光标所在处
	nnoremap gf <Cmd>call VSCodeNotify('workbench.action.navigateForward')<CR>
	" 跳转回上一次光标所在处
	nnoremap gb <Cmd>call VSCodeNotify('workbench.action.navigateBack')<CR>
	" 在当前分组内，切换vscode的editors
	nnoremap gy <Cmd>call VSCodeNotify('workbench.action.quickOpenPreviousRecentlyUsedEditorInGroup')<CR>
	" 切换tab
	" nnoremap gt <Cmd>call VSCodeNotify('workbench.action.openNextRecentlyUsedEditorInGroup')<CR>
	" 切换行注释
	nnoremap gc <Cmd>call VSCodeNotify('editor.action.commentLine')<CR>
	" 切换块注释
	nnoremap gt <Cmd>call VSCodeNotify('editor.action.blockComment')<CR>
	" 打开/关闭 codelens
	nnoremap gl <Cmd>call VSCodeNotify('codelens.showLensesInCurrentLine')<CR>
	" 展开所有折叠
	nnoremap zu <Cmd>call VSCodeNotify('editor.unfoldAll')<CR>
	" 关闭所有折叠
	nnoremap za <Cmd>call VSCodeNotify('editor.foldAll')<CR>
	" 展开当下折叠
	nnoremap zo <Cmd>call VSCodeNotify('editor.unfold')<CR>
	" 关闭当下折叠
	nnoremap zc <Cmd>call VSCodeNotify('editor.fold')<CR>
	" 切换当下折叠
	nnoremap zt <Cmd>call VSCodeNotify('editor.toggleFold')<CR>
	" 折叠所有注释
	nnoremap zm <Cmd>call VSCodeNotify('editor.foldAllBlockComments')<CR>
	" 转到文件中上一个问题
	nnoremap g[ <Cmd>call VSCodeNotify('editor.action.marker.prevInFiles')<CR>
	" 转到文件中下一个问题
	nnoremap g] <Cmd>call VSCodeNotify('editor.action.marker.nextInFiles')<CR>
endif
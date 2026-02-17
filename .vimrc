" Local Vim config for ptb (loaded via :set exrc + :set nosecure)

if exists("g:ptb_local_loaded")
  finish
endif
let g:ptb_local_loaded = 1

function! s:Slugify(title) abort
  let s = tolower(a:title)
  let s = substitute(s, '[^a-z0-9[:space:]-]', '', 'g')
  let s = substitute(s, '\s\+', '-', 'g')
  let s = substitute(s, '-\+', '-', 'g')
  let s = substitute(s, '^-', '', '')
  let s = substitute(s, '-$', '', '')
  return s
endfunction

function! s:NewPost(title) abort
  let slug = s:Slugify(a:title)
  if empty(slug)
    echoerr "Title produced an empty slug."
    return
  endif

  let filename = strftime('%Y%m%d') . '_' . slug . '.txt'
  let path = 'txt/' . filename

  execute 'edit ' . fnameescape(path)
  if line('$') == 1 && getline(1) ==# ''
    call setline(1, toupper(a:title))
    call append(1, '')
    normal! G
  endif
endfunction

" Usage: :NewPost unix as an ide
command! -nargs=+ NewPost call s:NewPost(<q-args>)

augroup ptb_txt_html_helpers
  autocmd!
  autocmd BufEnter,BufNewFile txt/*.txt setlocal textwidth=80

  " Insert HTML quickly while writing posts.
  autocmd BufEnter,BufNewFile txt/*.txt inoremap <buffer> ,a <a href=""></a><Left><Left><Left><Left>
  autocmd BufEnter,BufNewFile txt/*.txt inoremap <buffer> ,p <p></p><Left><Left><Left><Left>
  autocmd BufEnter,BufNewFile txt/*.txt inoremap <buffer> ,b <b></b><Left><Left><Left><Left>
  autocmd BufEnter,BufNewFile txt/*.txt inoremap <buffer> ,i <i></i><Left><Left><Left><Left>
  autocmd BufEnter,BufNewFile txt/*.txt inoremap <buffer> ,pre <pre></pre><Left><Left><Left><Left><Left><Left>
augroup END

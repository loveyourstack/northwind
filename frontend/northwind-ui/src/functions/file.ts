import ax from '@/api'

// requires content-disposition header filename to be set, e.g. "attachment; filename=Countries.xlsx"
// adapted from https://stackoverflow.com/questions/53772331/vue-html-js-how-to-download-a-file-to-browser-using-the-download-tag
export function fileDownload(url: string) {
    ax.get(url, { responseType: 'blob' })
      .then(resp => {

        // exit if json response type
        if (resp.data.type == 'application/json') {
          return
        }

        // get attachment filename from content-disposition header if possible
        var fileName = 'unknown'
        var contDispVal: string = resp.headers['content-disposition']
        //console.log('contDispVal: ' + contDispVal)
        if (contDispVal) {
            var cdFileName: string = contDispVal.split('filename=')[1]!
            //console.log('cdFileName: ' + cdFileName)
            if (cdFileName) {
                fileName = cdFileName
            }
        }

        const blob = new Blob([resp.data], { type: 'application/octet-stream' })
        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = fileName

        link.click()
        URL.revokeObjectURL(link.href)
      })
      .catch() // handled by interceptor
  }
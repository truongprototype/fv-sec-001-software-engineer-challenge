> Đọc toàn bộ README.md của dự án.
>
> Chưa viết code và chưa tự tạo dữ liệu kết quả.
>
> Hãy thực hiện:
>
> 1. Trích xuất tất cả yêu cầu chức năng, kỹ thuật và bàn giao.
> 2. Liệt kê các trường hợp biên và điểm chưa được đặc tả rõ.
> 3. Đề xuất kiến trúc Go tối giản, dễ kiểm thử và xử lý CSV khoảng 1 GB bằng streaming.
> 4. Xác định thuật toán tổng hợp và tìm Top 10.
> 5. Đề xuất cách xử lý số tiền, chia cho 0, dữ liệu sai định dạng và các trường hợp đồng hạng.
> 6. Đưa ra kế hoạch test và benchmark có thể tái lập.
>
> Tạo IMPLEMENTATION_PLAN.md. Phân biệt yêu cầu bắt buộc với đề xuất của bạn. Không tự ý bổ sung tính năng ngoài phạm vi.

> Các đề xuất trong  IMPLEMENTATION_PLAN.md đã được duyệt, hãy tạo bộ khung dự án Go.
> Yêu cầu:
>
> - Dùng Go standard library nếu chưa có lý do rõ ràng để thêm dependency.
> - Tách CLI, logic aggregation và phần đọc/ghi CSV thành các package hoặc module hợp lý.
> - Có thể chạy bằng go run và build thành binary bằng go build.
> - Hỗ trợ --input và --output.
> - Xác thực tham số, xử lý lỗi và trả về exit code phù hợp.
> - Không triển khai toàn bộ thuật toán trong một file main.go.
> - Thêm .gitignore và hướng dẫn chạy cơ bản.
> Chỉ triển khai bộ khung, có thể chưa cần tối ưu hiệu năng. Chạy gofmt, go vet và go test ./...; báo rõ lệnh nào đã chạy và kết quả thực tế.

> Triển khai riêng phần "3A — CSV reader và aggregator" theo IMPLEMENTATION_PLAN.md.
>
> Điều kiện:
>
> - Chỉ sửa những file cần thiết cho phần này.
> - Viết unit test cho logic mới và các trường hợp biên.
> - Không thay đổi hành vi đã được kiểm thử ở các phần trước nếu không có lý do.
> - Không tối ưu vi mô khi chưa có benchmark.
> - Không thêm dependency nếu thư viện chuẩn đáp ứng yêu cầu.
> - Chạy gofmt, go vet và go test ./... sau khi hoàn thành.
> - Báo cáo file đã sửa, quyết định kỹ thuật, kết quả test và các vấn đề còn tồn tại.
>
> Chỉ báo cáo test thành công nếu thực sự chạy đúng.

> Triển khai riêng phần "3B — Metrics và ranking" theo IMPLEMENTATION_PLAN.md.
>
> Điều kiện:
>
> - Chỉ sửa những file cần thiết cho phần này.
> - Viết unit test cho logic mới và các trường hợp biên.
> - Không thay đổi hành vi đã được kiểm thử ở các phần trước nếu không có lý do.
> - Không tối ưu vi mô khi chưa có benchmark.
> - Không thêm dependency nếu thư viện chuẩn đáp ứng yêu cầu.
> - Chạy gofmt, go vet và go test ./... sau khi hoàn thành.
> - Báo cáo file đã sửa, quyết định kỹ thuật, kết quả test và các vấn đề còn tồn tại.
> Chỉ báo cáo test thành công nếu thực sự chạy đúng.

> Triển khai riêng phần "3C — CSV writer và CLI integration" theo IMPLEMENTATION_PLAN.md.
>
> Điều kiện:
>
> - Chỉ sửa những file cần thiết cho phần này.
> - Viết unit test cho logic mới và các trường hợp biên.
> - Không thay đổi hành vi đã được kiểm thử ở các phần trước nếu không có lý do.
> - Không tối ưu vi mô khi chưa có benchmark.
> - Không thêm dependency nếu thư viện chuẩn đáp ứng yêu cầu.
> - Chạy gofmt, go vet và go test ./... sau khi hoàn thành.
> - Báo cáo file đã sửa, quyết định kỹ thuật, kết quả test và các vấn đề còn tồn tại.
> Chỉ báo cáo test thành công nếu thực sự chạy đúng.

> Hãy review repository như một reviewer độc lập. Đọc README.md và IMPLEMENTATION_PLAN.md, sau đó kiểm tra code và test hiện có.
>
> Ưu tiên tìm:
>
> 1. Sai logic aggregation, CTR, CPA hoặc ranking.
> 2. Lỗi khi impressions hoặc conversions bằng 0.
> 3. Sai lệch độ chính xác khi cộng spend.
> 4. Lỗi xử lý header, quoted CSV fields, dòng lỗi và file rỗng.
> 5. Bộ nhớ tăng theo số dòng thay vì số campaign.
> 6. Tài nguyên không được đóng, file output bị ghi dở hoặc lỗi bị nuốt.
> 7. Test thiếu hoặc test không thực sự kiểm chứng kết quả.
>
> Chưa sửa code. Trả về từng phát hiện gồm mức độ nghiêm trọng, vị trí, cách tái hiện và đề xuất sửa. Phân biệt lỗi đã xác minh với rủi ro giả định.

> - Spend và CPA được xuất với hai chữ số thập phân. Nếu input hợp lệ có độ chính xác vượt quá cent, phần hiển thị bị làm tròn; README cũng không nói không được tính trường hợp này, vậy thì hãy triển khai thêm.
> - Những giả định khác chưa cần làm.

> Chuẩn bị benchmark cho ứng dụng Ad Performance Aggregator.
>
> Yêu cầu:
>
> 1. Xác minh kết quả bằng fixture nhỏ có expected output đã biết.
> 2. Benchmark file thực tế (ad_data.csv); ghi rõ kích thước file, số dòng, số campaign, môi trường và phiên bản Go.
> 3. Đo wall-clock time và peak memory bằng công cụ phù hợp.
> 4. Phân biệt thời gian xử lý ứng dụng với thời gian tạo dữ liệu hoặc giải nén.
> 5. Lưu lệnh chạy, stdout/stderr và kết quả đo vào thư mục benchmarks/.
> 6. Nếu hiệu năng chưa tốt, phân tích nguyên nhân trước khi thay đổi code.
> 7. Sau mỗi tối ưu, chạy lại toàn bộ test và so sánh kết quả đầu ra với bản trước.
>
> Không bịa số liệu benchmark. Nếu chưa chạy được trên file thực, ghi rõ benchmark chỉ là thử nghiệm trên dữ liệu nhỏ.

> Hãy bổ sung khả năng thu thập CPU profile cho ứng dụng Go Ad Performance Aggregator để xác định nguyên nhân tiêu thụ CPU.
>
> Bối cảnh:
>
> - Ứng dụng xử lý CSV khoảng 1 GB theo streaming.
> - Có khoảng 50 campaign.
> - User CPU time gần bằng wall-clock time.
> - Hiện tại chưa có CPU profile và chưa được phép tối ưu code.
>
> Yêu cầu:
>
> 1. Kiểm tra cấu trúc dự án và cách chạy CLI hiện tại trước khi thay đổi.
> 2. Ưu tiên sử dụng runtime/pprof và go tool pprof của Go standard toolchain.
> 3. Đề xuất cách thu CPU profile trong một lần chạy benchmark trên file thực tế. Ưu tiên phương án ít ảnh hưởng đến code và không thay đổi logic xử lý.
> 4. Nếu cần sửa code, hãy thêm cơ chế profiling tùy chọn, chẳng hạn cờ --cpuprofile, chỉ bật khi người dùng yêu cầu. Phải start và stop profile đúng cách, đóng file an toàn và xử lý lỗi.
> 5. Không bật profiling mặc định. Không thay đổi thuật toán aggregation, CSV parsing, kiểu số hay luồng xử lý.
> 6. Chạy gofmt, go test ./... và go vet ./... nếu môi trường cho phép.
> 7. Hướng dẫn chính xác cách build, chạy benchmark với profile, xem top CPU functions và tạo call graph hoặc flame graph nếu công cụ hỗ trợ.
> 8. Không tự kết luận nguyên nhân chậm trước khi có profile thực tế.
>
> Hãy báo cáo kế hoạch trước, sau đó mới thực hiện các thay đổi tối thiểu cần thiết.


> Hoàn thiện repository để bàn giao.
>
> Yêu cầu:
>
> - Xác nhận go test ./..., go vet ./... và build đều thành công.
> - Tạo binary cho hệ điều hành và kiến trúc đích đã xác định.
> - Kiểm tra binary thực sự chạy được với --input và --output.
> - Làm lại README.md sạch: cài đặt, cách chạy, thiết kế, thư viện, benchmark và peak memory.
> - Đảm bảo top10_ctr.csv và top10_cpa.csv được tạo từ dữ liệu thực, không phải dữ liệu minh họa.
> - Kiểm tra các file bắt buộc và chất lượng Git diff.
>
> Không đưa số liệu chưa đo hoặc kết quả chưa xác minh vào tài liệu. Kết thúc bằng checklist bàn giao và các giới hạn còn tồn tại.

> Hãy bổ sung Dockerfile cho dự án Go Ad Performance Aggregator hiện tại và thực sự kiểm tra khả năng build, chạy ứng dụng trong Docker.
>
> ## 1. Tạo Dockerfile
> - Ưu tiên multi-stage build: dùng Go image để biên dịch và một runtime image nhỏ để chạy binary.
> - Ưu tiên Docker official images và phiên bản Go phù hợp với go.mod của dự án.
> - Build binary phù hợp với môi trường container Linux.
> - Nếu không cần thiết, không đưa source code, Git history hoặc dữ liệu CSV vào runtime image.
> - Sử dụng .dockerignore để loại trừ các file không cần thiết.
> - Không chạy ứng dụng bằng root nếu có thể cấu hình non-root user mà không làm hỏng chức năng.
> - Thiết kế để file input và thư mục output được truyền qua bind mount hoặc volume.
> - Không copy dữ liệu quảng cáo dung lượng khoảng 1 GB vào image.
>
> ## 2. Build image
> Chọn tên image nhất quán. Nếu build thất bại, tìm nguyên nhân, sửa trong phạm vi cần thiết rồi build lại. Không tuyên bố thành công nếu chưa chạy lệnh và xác nhận kết quả.
>
> ## 3. Chạy thử trong Docker
> - Tạo hoặc sử dụng fixture CSV nhỏ có expected output được xác định độc lập.
> - Chạy container với file input được mount ở chế độ chỉ đọc và thư mục output được mount để lấy kết quả ra host.
> - Kiểm tra exit code, log, các file CSV đầu ra, header, số dòng và nội dung so với expected output.
> - Kiểm tra ít nhất một trường hợp lỗi phù hợp, ví dụ input không tồn tại hoặc tham số không hợp lệ.
> - Không dùng file CSV giả để tuyên bố đã kiểm thử hiệu năng trên dữ liệu 1 GB.
>
> ## 4. Kiểm tra chất lượng
> - Chạy go test ./..., go vet ./... và gofmt nếu có thay đổi code Go.
> - Xác minh Docker build có thể tái lập theo hướng dẫn trong README.
> - Kiểm tra container không phụ thuộc vào file hoặc đường dẫn chỉ tồn tại trên máy thật.
> - Không sửa thuật toán aggregation hay tối ưu hiệu năng ngoài phạm vi yêu cầu.
>
> ## 5. Cập nhật tài liệu
> Cập nhật README.md với:
> - Thêm Docker.
> - Lệnh build image.
> - Lệnh chạy container với input/output mount.
> - Cách chạy test.
> - Các giới hạn hoặc lưu ý về quyền truy cập thư mục và file.
>
> ## 6. Báo cáo cuối
> Sau khi thực hiện, hãy báo cáo:
> - File đã tạo hoặc sửa.
> - Docker image và chiến lược build.
> - Các lệnh thực tế đã chạy cùng kết quả.
> - Kết quả kiểm tra output với expected output.
> - Các test đã chạy và kết quả.
> - Những việc chưa thể kiểm chứng hoặc còn tồn tại.

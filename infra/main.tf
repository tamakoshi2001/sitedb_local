provider "aws" {
  region = "ap-northeast-1" # 必要に応じてこのリージョンを変更してください
}

variable "OPENAI_API_KEY" {
  description = "API Key for the application"
  type        = string
}

variable "BUCKET_NAME" {
  description = "API Key for the application"
  type        = string
}

variable "AWS_ACCESS_KEY_ID" {
    description = "API Key for AWS"
    type = string
}

variable "AWS_SECRET_ACCESS_KEY" {
    description = "Secret Key for AWS"
    type = string
}

resource "aws_key_pair" "example_keypair" {
  key_name   = "example_keypair"
  public_key = file("~/.ssh/id_rsa.pub") # こちらのパスは、公開鍵の実際のパスに置き換えてください。
}

resource "aws_s3_bucket" "example_bucket" {
  bucket = "${var.BUCKET_NAME}" # バケット名を必要に応じて変更してください
}

resource "aws_instance" "example_instance" {
  ami           = "ami-0fd8f5842685ca887" # Amazon Linux 2 LTS のAMI ID。リージョンや最新のAMIに応じて更新が必要
  instance_type = "t2.micro"
  key_name      = aws_key_pair.example_keypair.key_name
  security_groups = [aws_security_group.example_sg.name]  # セキュリティグループを関連付け

  user_data = <<-EOF
              #!/bin/bash
              echo "BUCKET_NAME=${aws_s3_bucket.example_bucket.bucket}" >> /etc/environment
              echo "OPENAI_API_KEY=${var.OPENAI_API_KEY}" >> /etc/environment
              echo "AWS_ACCESS_KEY_ID=${var.AWS_ACCESS_KEY_ID}" >> /etc/environment
              echo "AWS_SECRET_ACCESS_KEY=${var.AWS_SECRET_ACCESS_KEY}" >> /etc/environment

              sudo yum update -y
              sudo yum install git -y

              sudo yum install wget -y
              sudo -u ec2-user wget -O /home/ec2-user/go1.21.0.linux-amd64.tar.gz https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
              sudo tar -C /usr/local -xzf /home/ec2-user/go1.21.0.linux-amd64.tar.gz
              echo "export PATH=$PATH:/usr/local/go/bin" > /etc/profile.d/go.sh
              source /etc/profile.d/go.sh
              cd /home/ec2-user
              git clone https://github.com/tamakoshi2001/sitedb.git
              cd sitedb
              sudo chown -R ec2-user:ec2-user /home/ec2-user/sitedb
              sudo -u ec2-user /usr/local/go/bin/go run main.go
              EOF

  tags = {
    Name = "example-instance"
  }
}

resource "aws_security_group" "example_sg" {
  name        = "example_security_group1"
  description = "Allow inbound traffic on port 8080"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # このCIDRはインターネット全体からのアクセスを許可します。適切な範囲に制限することを検討してください。
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "example_sg"
  }
}

output "bucket_arn" {
  value = aws_s3_bucket.example_bucket.arn
}

output "instance_public_ip" {
  value = aws_instance.example_instance.public_ip
}

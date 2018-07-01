namespace ReceiptEditor
{
    partial class ReceiptEditor
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.splitContainer1 = new System.Windows.Forms.SplitContainer();
            this.imageBox = new System.Windows.Forms.PictureBox();
            this.imageNameBox = new System.Windows.Forms.TextBox();
            this.imageFolderBox = new System.Windows.Forms.TextBox();
            this.nextButton = new System.Windows.Forms.Button();
            this.scanResultsTextBox = new System.Windows.Forms.TextBox();
            this.processedFolderBox = new System.Windows.Forms.TextBox();
            this.label2 = new System.Windows.Forms.Label();
            this.receiptFolderBox = new System.Windows.Forms.TextBox();
            this.scanButton = new System.Windows.Forms.Button();
            this.label1 = new System.Windows.Forms.Label();
            ((System.ComponentModel.ISupportInitialize)(this.splitContainer1)).BeginInit();
            this.splitContainer1.Panel1.SuspendLayout();
            this.splitContainer1.Panel2.SuspendLayout();
            this.splitContainer1.SuspendLayout();
            ((System.ComponentModel.ISupportInitialize)(this.imageBox)).BeginInit();
            this.SuspendLayout();
            // 
            // splitContainer1
            // 
            this.splitContainer1.Dock = System.Windows.Forms.DockStyle.Fill;
            this.splitContainer1.Location = new System.Drawing.Point(0, 0);
            this.splitContainer1.Name = "splitContainer1";
            // 
            // splitContainer1.Panel1
            // 
            this.splitContainer1.Panel1.Controls.Add(this.imageBox);
            // 
            // splitContainer1.Panel2
            // 
            this.splitContainer1.Panel2.Controls.Add(this.imageNameBox);
            this.splitContainer1.Panel2.Controls.Add(this.imageFolderBox);
            this.splitContainer1.Panel2.Controls.Add(this.nextButton);
            this.splitContainer1.Panel2.Controls.Add(this.scanResultsTextBox);
            this.splitContainer1.Panel2.Controls.Add(this.processedFolderBox);
            this.splitContainer1.Panel2.Controls.Add(this.label2);
            this.splitContainer1.Panel2.Controls.Add(this.receiptFolderBox);
            this.splitContainer1.Panel2.Controls.Add(this.scanButton);
            this.splitContainer1.Panel2.Controls.Add(this.label1);
            this.splitContainer1.Size = new System.Drawing.Size(874, 428);
            this.splitContainer1.SplitterDistance = 604;
            this.splitContainer1.TabIndex = 0;
            // 
            // imageBox
            // 
            this.imageBox.Dock = System.Windows.Forms.DockStyle.Fill;
            this.imageBox.Location = new System.Drawing.Point(0, 0);
            this.imageBox.Name = "imageBox";
            this.imageBox.Size = new System.Drawing.Size(604, 428);
            this.imageBox.SizeMode = System.Windows.Forms.PictureBoxSizeMode.Zoom;
            this.imageBox.TabIndex = 0;
            this.imageBox.TabStop = false;
            this.imageBox.Paint += new System.Windows.Forms.PaintEventHandler(this.imageBox_Paint);
            this.imageBox.MouseDown += new System.Windows.Forms.MouseEventHandler(this.imageBox_MouseDown);
            this.imageBox.MouseMove += new System.Windows.Forms.MouseEventHandler(this.imageBox_MouseMove);
            this.imageBox.MouseUp += new System.Windows.Forms.MouseEventHandler(this.imageBox_MouseUp);
            // 
            // imageNameBox
            // 
            this.imageNameBox.Location = new System.Drawing.Point(2, 367);
            this.imageNameBox.Name = "imageNameBox";
            this.imageNameBox.ReadOnly = true;
            this.imageNameBox.Size = new System.Drawing.Size(248, 20);
            this.imageNameBox.TabIndex = 8;
            // 
            // imageFolderBox
            // 
            this.imageFolderBox.Location = new System.Drawing.Point(2, 341);
            this.imageFolderBox.Name = "imageFolderBox";
            this.imageFolderBox.ReadOnly = true;
            this.imageFolderBox.Size = new System.Drawing.Size(248, 20);
            this.imageFolderBox.TabIndex = 7;
            // 
            // nextButton
            // 
            this.nextButton.Location = new System.Drawing.Point(179, 393);
            this.nextButton.Name = "nextButton";
            this.nextButton.Size = new System.Drawing.Size(75, 23);
            this.nextButton.TabIndex = 6;
            this.nextButton.Text = "Next";
            this.nextButton.UseVisualStyleBackColor = true;
            this.nextButton.Click += new System.EventHandler(this.nextButton_Click);
            // 
            // scanResultsTextBox
            // 
            this.scanResultsTextBox.Location = new System.Drawing.Point(6, 87);
            this.scanResultsTextBox.Name = "scanResultsTextBox";
            this.scanResultsTextBox.ReadOnly = true;
            this.scanResultsTextBox.Size = new System.Drawing.Size(248, 20);
            this.scanResultsTextBox.TabIndex = 5;
            // 
            // processedFolderBox
            // 
            this.processedFolderBox.Location = new System.Drawing.Point(85, 32);
            this.processedFolderBox.Name = "processedFolderBox";
            this.processedFolderBox.Size = new System.Drawing.Size(169, 20);
            this.processedFolderBox.TabIndex = 4;
            this.processedFolderBox.Text = "C:\\Users\\Gustave\\Desktop\\Data Archive\\processed_receipts";
            // 
            // label2
            // 
            this.label2.AutoSize = true;
            this.label2.Location = new System.Drawing.Point(3, 32);
            this.label2.Name = "label2";
            this.label2.Size = new System.Drawing.Size(76, 13);
            this.label2.TabIndex = 3;
            this.label2.Text = "Processed Dir.";
            // 
            // receiptFolderBox
            // 
            this.receiptFolderBox.Location = new System.Drawing.Point(85, 6);
            this.receiptFolderBox.Name = "receiptFolderBox";
            this.receiptFolderBox.Size = new System.Drawing.Size(169, 20);
            this.receiptFolderBox.TabIndex = 2;
            this.receiptFolderBox.Text = "C:\\Users\\Gustave\\Desktop\\Data Archive\\receipts";
            // 
            // scanButton
            // 
            this.scanButton.Location = new System.Drawing.Point(179, 58);
            this.scanButton.Name = "scanButton";
            this.scanButton.Size = new System.Drawing.Size(75, 23);
            this.scanButton.TabIndex = 1;
            this.scanButton.Text = "Scan";
            this.scanButton.UseVisualStyleBackColor = true;
            this.scanButton.Click += new System.EventHandler(this.scanButton_Click);
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(3, 9);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(76, 13);
            this.label1.TabIndex = 0;
            this.label1.Text = "Receipt Folder";
            // 
            // ReceiptEditor
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(874, 428);
            this.Controls.Add(this.splitContainer1);
            this.Name = "ReceiptEditor";
            this.Text = "Receipt Editor";
            this.Resize += new System.EventHandler(this.ReceiptEditor_Resize);
            this.splitContainer1.Panel1.ResumeLayout(false);
            this.splitContainer1.Panel2.ResumeLayout(false);
            this.splitContainer1.Panel2.PerformLayout();
            ((System.ComponentModel.ISupportInitialize)(this.splitContainer1)).EndInit();
            this.splitContainer1.ResumeLayout(false);
            ((System.ComponentModel.ISupportInitialize)(this.imageBox)).EndInit();
            this.ResumeLayout(false);

        }

        #endregion

        private System.Windows.Forms.SplitContainer splitContainer1;
        private System.Windows.Forms.PictureBox imageBox;
        private System.Windows.Forms.TextBox imageNameBox;
        private System.Windows.Forms.TextBox imageFolderBox;
        private System.Windows.Forms.Button nextButton;
        private System.Windows.Forms.TextBox scanResultsTextBox;
        private System.Windows.Forms.TextBox processedFolderBox;
        private System.Windows.Forms.Label label2;
        private System.Windows.Forms.TextBox receiptFolderBox;
        private System.Windows.Forms.Button scanButton;
        private System.Windows.Forms.Label label1;
    }
}

